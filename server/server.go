package server

import (
	"github.com/HistoryLabs/events-api/server/routes"
	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()

	router.Use(CorsMiddleware)

	router.GET("/", routes.Home)
	router.GET("/date", routes.FetchDate)
	router.GET("/year/*year", routes.FetchYear)

	router.Run("localhost:5000")
}

func CorsMiddleware(c *gin.Context) {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Next()
}
