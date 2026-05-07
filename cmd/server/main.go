package main

import (
	"log"

	"github.com/HistoryLabs/events-api/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	router := gin.Default()
	router.Use(corsMiddleware)
	router.GET("/", routes.Home)
	router.GET("/date", routes.FetchDate)
	router.GET("/year/*year", routes.FetchYear)
	router.Run("localhost:5000")
}

func corsMiddleware(c *gin.Context) {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Next()
}
