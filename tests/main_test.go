package tests

import (
	"github.com/dencat/fixss/fixss"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	router := fixss.CreateRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/quoteConfig", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{}\n", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/quoteConfig", strings.NewReader("a"))
	router.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"status\":\"error\"}\n", w.Body.String())

	fixss.LoadDefaultQuoteConfig()
	content, err := ioutil.ReadFile("get_quote_config_expected_1.json")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/quoteConfig", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, removeAllWhiteSpace(string(content)), removeAllWhiteSpace(w.Body.String()))

	content, err = ioutil.ReadFile("set_quotes.json")
	if err != nil {
		assert.Fail(t, err.Error())
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/quoteConfig", strings.NewReader(string(content)))
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"status\":\"ok\"}\n", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/quoteConfig", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"EUR/USD_TOM\":{\"symbol\":\"EUR/USD_TOM\",\"interval\":5000,\"entities\":[{\"size\":1000,\"direction\":\"bid\",\"price\":1.03},{\"size\":1000,\"direction\":\"offer\",\"price\":1.25},{\"size\":1000000,\"direction\":\"bid\",\"price\":1.05}]}}\n", w.Body.String())

	assert.Equal(t, true, fixss.GetQuoteConfig("EUR/USD_TOD") == nil)
	assert.Equal(t, true, fixss.GetQuoteConfig("EUR/USD_TOM") != nil)
	assert.Equal(t, int64(5000), fixss.GetQuoteConfig("EUR/USD_TOM").Interval)
	assert.Equal(t, 3, len(fixss.GetQuoteConfig("EUR/USD_TOM").Entities))
}

func removeAllWhiteSpace(str string) string {
	return strings.Replace(strings.Replace(string(str), "\n", "", -1), " ", "", -1)
}
