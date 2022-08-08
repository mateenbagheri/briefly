package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/controllers"
)

func UrlRoute(router *gin.Engine) {
	url := router.Group("/url")
	{
		url.POST("/", controllers.CreateURL)
		url.GET("/:ShortenedUrl", controllers.GetURLByShortened)

		url.GET("/user/:UserID")
		url.GET("/collection/:CollectionID")

		url.DELETE("/")
	}
}
