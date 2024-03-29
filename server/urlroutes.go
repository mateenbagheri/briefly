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
		url.GET("/collection/:CollectionID", controllers.GetCollectionURLs)
		url.GET("/user/:UserID", controllers.GetUserURLs)
		url.GET("/user/:UserID/history", controllers.GetUserUrlHistory)
		url.GET("/user/:UserID/report", controllers.GetUserUrlReport)

		url.DELETE("/:LinkID", controllers.DeleteURLByID)
	}
}
