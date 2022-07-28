package server

import "github.com/gin-gonic/gin"

func CollectionRoute(router *gin.Engine) {
	collection := router.Group("/collection")
	{
		collection.POST("/")
		collection.GET("/")
		collection.PUT("/")
		collection.DELETE("/")
	}
}
