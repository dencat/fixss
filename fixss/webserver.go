package fixss

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	log "github.com/jeanphorn/log4go"
	"net/http"
	"strconv"
)

const API_PATH = "/api/v1/"

var statusOk = gin.H{"status": "ok"}
var statusError = gin.H{"status": "error"}

func CreateRouter() *gin.Engine {
	router := gin.Default()

	router.GET(API_PATH+"quoteConfig", getQuoteConfig)
	router.POST(API_PATH+"quoteConfig", setQuoteConfig)
	router.POST(API_PATH+"quoteConfigs", setQuoteConfigs)

	router.GET(API_PATH+"orderConfig", getOrderConfig)
	router.POST(API_PATH+"orderConfig", setOrderConfig)

	pprof.Register(router)

	return router
}

func StartWebServer(port int) {
	go func() {
		router := CreateRouter()
		router.Run(":" + strconv.Itoa(port))
	}()
}

func getQuoteConfig(context *gin.Context) {
	context.JSON(http.StatusOK, quoteConfigs)
}

func setQuoteConfig(context *gin.Context) {
	var quoteConfig QuoteConfig
	err := context.BindJSON(&quoteConfig)

	if err != nil {
		log.Error("Set quote config error: %s", err.Error())
		context.AbortWithStatusJSON(http.StatusInternalServerError, statusError)
		return
	}
	SetQuoteConfig(quoteConfig)
	context.JSON(http.StatusOK, statusOk)
}

func setQuoteConfigs(context *gin.Context) {
	var quoteConfigs []QuoteConfig
	err := context.BindJSON(&quoteConfigs)

	if err != nil {
		log.Error("Set quote config error: %s", err.Error())
		context.AbortWithStatusJSON(http.StatusInternalServerError, statusError)
		return
	}
	for _, quoteConfig := range quoteConfigs {
		SetQuoteConfig(quoteConfig)
	}

	context.JSON(http.StatusOK, statusOk)
}

func getOrderConfig(context *gin.Context) {
	context.JSON(http.StatusOK, orderConfigs)
}

func setOrderConfig(context *gin.Context) {
	var orderConfig OrderConfig
	err := context.BindJSON(&orderConfig)

	if err != nil {
		log.Error("Set order config error: %s", err.Error())
		context.AbortWithStatusJSON(http.StatusInternalServerError, statusError)
		return
	}
	SetOrderConfig(orderConfig)
	context.JSON(http.StatusOK, statusOk)
}
