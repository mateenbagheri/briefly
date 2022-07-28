package server

import "github.com/gin-gonic/gin"

func UrlRoute(router *gin.Engine) {
	url := router.Group("/url")
	{
		url.GET("")
		url.POST("")
		url.PUT("")
		url.DELETE("")
	}
}
