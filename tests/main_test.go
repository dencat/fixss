package tests

import (
	"github.com/dencat/fixss/fixss"
	"github.com/gin-gonic/gin"
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

	config := fixss.Config{}
	config.Fix.Config = "./config/server.cfg"
	config.Quote.Config = "./config/quoteDefaultConfig.json"
	err := fixss.StartAcceptor(&config)
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
	asrt.Equal("{}", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/quoteConfig", strings.NewReader("a"))
	router.ServeHTTP(w, req)
	asrt.Equal(400, w.Code)
	asrt.Equal("{\"status\":\"error\"}", w.Body.String())

	err = fixss.LoadDefaultQuoteConfig(&config)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	content, err := ioutil.ReadFile("get_quote_config_expected_1.json")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/quoteConfig", nil)
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal(removeAllWhiteSpace(string(content)), removeAllWhiteSpace(w.Body.String()))

	sendPostFromFile("/api/v1/quoteConfig", "set_quote.json", t, router, asrt)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/quoteConfig", nil)
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal("{\"EUR/USD_TOM\":{\"symbol\":\"EUR/USD_TOM\",\"interval\":100,\"entities\":[{\"size\":1000,\"direction\":\"bid\",\"minPrice\":1.031,\"maxPrice\":1.031},{\"size\":1000,\"direction\":\"offer\",\"minPrice\":1.251,\"maxPrice\":1.251},{\"size\":1000000,\"direction\":\"bid\",\"minPrice\":1.051,\"maxPrice\":1.051}]}}", w.Body.String())

	asrt.Equal(true, fixss.GetQuoteConfig("EUR/USD_TOD") == nil)
	asrt.Equal(true, fixss.GetQuoteConfig("EUR/USD_TOM") != nil)
	asrt.Equal(int64(100), fixss.GetQuoteConfig("EUR/USD_TOM").Interval)
	asrt.Equal(3, len(fixss.GetQuoteConfig("EUR/USD_TOM").Entities))

	clientApp.SendMarketDataRequest("EUR/USD_TOM")

	for {
		if clientApp.GetLastQuote("EUR/USD_TOM_0_1000") == "1.031" &&
			clientApp.GetLastQuote("EUR/USD_TOM_1_1000") == "1.251" &&
			clientApp.GetLastQuote("EUR/USD_TOM_0_1000000") == "1.051" {
			break
		}
	}

	sendPostFromFile("/api/v1/quoteConfig", "set_empty_quotes.json", t, router, asrt)

	for {
		if clientApp.GetLastQuote("EUR/USD_TOM") == "drop" {
			break
		}
	}

	sendPostFromFile("/api/v1/quoteConfigs", "set_quotes.json", t, router, asrt)

	for {
		if clientApp.GetLastQuote("EUR/USD_TOM_0_1000") == "1.031" &&
			clientApp.GetLastQuote("EUR/USD_TOM_1_1000") == "1.251" &&
			clientApp.GetLastQuote("EUR/USD_TOM_0_1000000") == "1.051" {
			break
		}
	}

	sendPostFromFile("/api/v1/orderConfig", "set_order_config_reject.json", t, router, asrt)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/orderConfig", nil)
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal("{\"EUR/USD_TOM\":{\"symbol\":\"EUR/USD_TOM\",\"strategy\":\"reject\"}}", w.Body.String())

	asrt.Equal(fixss.Reject, fixss.GetOrderConfig("EUR/USD_TOM").Strategy)
	clientApp.SendOrder("111", "EUR/USD_TOM", decimal.NewFromInt(1500), decimal.NewFromFloat(1.05), "1")
	waitOrderStatus(clientApp, "111", enum.OrdStatus_REJECTED)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/orderConfig", strings.NewReader("{\"symbol\":\"EUR/USD_TOM\",\"strategy\":\"accept\"}"))
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal("{\"status\":\"ok\"}", w.Body.String())
	asrt.Equal(fixss.Accept, fixss.GetOrderConfig("EUR/USD_TOM").Strategy)

	clientApp.SendOrder("222", "EUR/USD_TOM", decimal.NewFromInt(10000000), decimal.NewFromFloat(1.3), enum.Side_BUY)
	waitOrderStatus(clientApp, "222", enum.OrdStatus_REJECTED)

	clientApp.SendOrder("333", "EUR/USD_TOM", decimal.NewFromInt(500), decimal.NewFromFloat(1.249), enum.Side_BUY)
	waitOrderStatus(clientApp, "333", enum.OrdStatus_REJECTED)

	clientApp.SendOrder("444", "EUR/USD_TOM", decimal.NewFromInt(500), decimal.NewFromFloat(1.251), enum.Side_BUY)
	waitOrderStatus(clientApp, "444", enum.OrdStatus_FILLED)

	clientApp.SendOrder("555", "EUR/USD_TOM", decimal.NewFromInt(500), decimal.NewFromFloat(1.26), enum.Side_BUY)
	waitOrderStatus(clientApp, "555", enum.OrdStatus_FILLED)

	clientApp.SendOrder("777", "EUR/USD_TOM", decimal.NewFromInt(500), decimal.NewFromFloat(1.032), enum.Side_SELL)
	waitOrderStatus(clientApp, "777", enum.OrdStatus_REJECTED)

	clientApp.SendOrder("888", "EUR/USD_TOM", decimal.NewFromInt(2500), decimal.NewFromFloat(1.032), enum.Side_SELL)
	waitOrderStatus(clientApp, "888", enum.OrdStatus_FILLED)

	client.Stop()
	fixss.StopAcceptor()
}

func waitOrderStatus(clientApp *TradeClient, orderId string, status enum.OrdStatus) {
	for {
		if clientApp.GetOrderStatus(orderId) == status {
			break
		}
	}
}

func sendPostFromFile(apiPath string, filePath string, t *testing.T, router *gin.Engine, asrt *assert.Assertions) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", apiPath, strings.NewReader(string(content)))
	router.ServeHTTP(w, req)
	asrt.Equal(200, w.Code)
	asrt.Equal("{\"status\":\"ok\"}", w.Body.String())
}

func removeAllWhiteSpace(str string) string {
	return strings.Replace(strings.Replace(string(str), "\n", "", -1), " ", "", -1)
}

func TestConfig(t *testing.T) {
	asrt := assert.New(t)
	cfg, err := fixss.NewConfig("config/config.yml")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	asrt.Equal("127.0.0.1", cfg.Server.Host)
	asrt.Equal("8888", cfg.Server.Port)
}
