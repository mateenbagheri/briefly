package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/controllers"
	"github.com/mateenbagheri/briefly/middleware"
)

func CollectionRoute(router *gin.Engine) {
	collection := router.Group("/collection").Use(middleware.RequireAuth)
	{
		collection.GET("/", controllers.GetAllCollections)

		collection.GET("/user/:UserID", controllers.GetUserCollections)
		collection.GET("/:CollectionID", controllers.GetCollectionByID)

		collection.DELETE("/:CollectionID", controllers.DeleteCollectionByID)

		collection.POST("/", controllers.CreateCollection)
		collection.POST("/url", controllers.AddUrlToCollection)

		collection.PUT("/:CollectionID", controllers.EditCollectionByID)
	}
}
