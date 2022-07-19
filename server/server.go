package server

import (
	"github.com/HistoryLabs/events-api/server/routes"
	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()

	router.GET("/date", routes.FetchDate)
	router.GET("/year", routes.FetchYear)

	router.Run("localhost:5000")
}
