package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mateenbagheri/briefly/server"
)

func main() {
	router := gin.Default()

	// adding up routings to our system.
	server.CollectionRoute(router)
	server.UrlRoute(router)
	server.UserRoute(router)

	// configs.ConnectDB()
	router.Run(":6000")
}
