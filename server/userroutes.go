package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/controllers"
	"github.com/mateenbagheri/briefly/middleware"
)

func UserRoute(router *gin.Engine) {
	user := router.Group("/user")
	{
		user.POST("/signup", controllers.SignUp)
		user.POST("/login", controllers.Login)

		user.GET("/validate", middleware.RequireAuth, controllers.Validate)
	}
}
