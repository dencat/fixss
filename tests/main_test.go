package tests

import (
	"github.com/dencat/fixss/fixss"
	"github.com/quickfixgo/enum"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWithTimeOut(t *testing.T) {
	timeout := time.After(10 * time.Second)
	done := make(chan bool)
	go func() {
		testServer(t)
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}
}

func testServer(t *testing.T) {
	asrt := assert.New(t)

	err := fixss.StartAcceptor()
	asrt.NoError(err)

	loginDone := make(chan bool, 1)

	client, clientApp, err := CreateInitiator(loginDone)
	asrt.NoError(err)
	err = client.Start()
	asrt.NoError(err)

	//	wait login
	<-loginDone

	router := fixss.CreateRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/quoteConfig", nil)
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal("{}\n", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/quoteConfig", strings.NewReader("a"))
	router.ServeHTTP(w, req)
	asrt.Equal(400, w.Code)
	asrt.Equal("{\"status\":\"error\"}\n", w.Body.String())

	fixss.LoadDefaultQuoteConfig()
	content, err := ioutil.ReadFile("get_quote_config_expected_1.json")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/quoteConfig", nil)
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal(removeAllWhiteSpace(string(content)), removeAllWhiteSpace(w.Body.String()))

	content, err = ioutil.ReadFile("set_quotes.json")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/quoteConfig", strings.NewReader(string(content)))
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal("{\"status\":\"ok\"}\n", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/quoteConfig", nil)
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal("{\"EUR/USD_TOM\":{\"symbol\":\"EUR/USD_TOM\",\"interval\":100,\"entities\":[{\"size\":1000,\"direction\":\"bid\",\"price\":1.03},{\"size\":1000,\"direction\":\"offer\",\"price\":1.25},{\"size\":1000000,\"direction\":\"bid\",\"price\":1.05}]}}\n", w.Body.String())

	asrt.Equal(true, fixss.GetQuoteConfig("EUR/USD_TOD") == nil)
	asrt.Equal(true, fixss.GetQuoteConfig("EUR/USD_TOM") != nil)
	asrt.Equal(int64(100), fixss.GetQuoteConfig("EUR/USD_TOM").Interval)
	asrt.Equal(3, len(fixss.GetQuoteConfig("EUR/USD_TOM").Entities))

	clientApp.SendMarketDataRequest("EUR/USD_TOM")

	for {
		if clientApp.GetLastQuote("EUR/USD_TOM_0_1000") == "1.03" &&
			clientApp.GetLastQuote("EUR/USD_TOM_1_1000") == "1.25" &&
			clientApp.GetLastQuote("EUR/USD_TOM_0_1000000") == "1.05" {
			break
		}
	}

	content, err = ioutil.ReadFile("set_order_config.json")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/orderConfig", strings.NewReader(string(content)))
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal("{\"status\":\"ok\"}\n", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/orderConfig", nil)
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal("{\"EUR/USD_TOM\":{\"symbol\":\"EUR/USD_TOM\",\"strategy\":\"reject\"}}\n", w.Body.String())

	asrt.Equal(fixss.Reject, fixss.GetOrderConfig("EUR/USD_TOM").Strategy)
	clientApp.SendOrder("123", "EUR/USD_TOM", decimal.NewFromInt(1500), decimal.NewFromFloat(1.05))

	for {
		if clientApp.GetOrderStatus("123") == enum.OrdStatus_REJECTED {
			break
		}
	}

	client.Stop()
	fixss.StopAcceptor()
}

func removeAllWhiteSpace(str string) string {
	return strings.Replace(strings.Replace(string(str), "\n", "", -1), " ", "", -1)
}
