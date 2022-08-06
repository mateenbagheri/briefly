package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/controllers"
)

func CollectionRoute(router *gin.Engine) {
	collection := router.Group("/collection")
	{
		collection.GET("/", controllers.GetAllCollections)
		collection.GET("/:CollectionID", controllers.GetCollectionByID)

		collection.DELETE("/:CollectionID", controllers.DeleteCollectionByID)

		collection.POST("/", controllers.CreateCollection)
		collection.PUT("/:CollectionID", controllers.EditCollectionByID)
	}
}
