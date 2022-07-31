package server

import "github.com/gin-gonic/gin"

func UserRoute(router *gin.Engine) {
	user := router.Group("/user")
	{
		user.GET("/")
		user.POST("/")
		user.PUT("/")
		user.DELETE("/")
	}
}
