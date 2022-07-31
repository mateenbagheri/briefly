package server

import "github.com/gin-gonic/gin"

func CollectionRoute(router *gin.Engine) {
	collection := router.Group("/collection")
	{
		collection.POST("/")

		collection.GET("/")
		collection.GET("/:id")

		collection.PUT("/")

		collection.DELETE("/")
	}
}
