package fixss

import (
	"github.com/gin-gonic/gin"
)

const API_PATH = "/api/v1/"

func CreateRouter() *gin.Engine {
	router := gin.Default()
	router.GET(API_PATH+"quoteConfig", getQuoteConfig)
	router.POST(API_PATH+"quoteConfig", setQuoteConfig)
	return router
}
func StartWebServer() {
	go func() {
		router := CreateRouter()
		router.Run(":8080")
	}()
}

func getQuoteConfig(context *gin.Context) {
	context.JSON(200, configs)
}

func setQuoteConfig(context *gin.Context) {
	var quoteConfig QuoteConfig
	err := context.BindJSON(&quoteConfig)

	if err != nil {
		println(err.Error())
		context.AbortWithStatusJSON(500, gin.H{"status": "error"})
		return
	}
	SetQuoteConfig(quoteConfig)
	context.JSON(200, gin.H{"status": "ok"})
}
