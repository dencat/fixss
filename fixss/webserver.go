package fixss

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const API_PATH = "/api/v1/"

var statusOk = gin.H{"status": "ok"}
var statusError = gin.H{"status": "error"}

func CreateRouter() *gin.Engine {
	router := gin.Default()

	router.GET(API_PATH+"quoteConfig", getQuoteConfig)
	router.POST(API_PATH+"quoteConfig", setQuoteConfig)

	router.GET(API_PATH+"orderConfig", getOrderConfig)
	router.POST(API_PATH+"orderConfig", setOrderConfig)

	return router
}

func StartWebServer() {
	go func() {
		router := CreateRouter()
		router.Run(":8080")
	}()
}

func getQuoteConfig(context *gin.Context) {
	context.JSON(http.StatusOK, quoteConfigs)
}

func setQuoteConfig(context *gin.Context) {
	var quoteConfig QuoteConfig
	err := context.BindJSON(&quoteConfig)

	if err != nil {
		Log.Errorf("Set quote config error: %s", err.Error())
		context.AbortWithStatusJSON(http.StatusInternalServerError, statusError)
		return
	}
	SetQuoteConfig(quoteConfig)
	context.JSON(http.StatusOK, statusOk)
}

func getOrderConfig(context *gin.Context) {
	context.JSON(http.StatusOK, orderConfigs)
}

func setOrderConfig(context *gin.Context) {
	var orderConfig OrderConfig
	err := context.BindJSON(&orderConfig)

	if err != nil {
		Log.Errorf("Set order config error: %s", err.Error())
		context.AbortWithStatusJSON(http.StatusInternalServerError, statusError)
		return
	}
	SetOrderConfig(orderConfig)
	context.JSON(http.StatusOK, statusOk)
}
