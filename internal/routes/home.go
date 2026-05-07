package routes

import "github.com/gin-gonic/gin"

func Home(c *gin.Context) {
	c.Data(200, "text/html", []byte("Welcome to the Events API! You can find the documentation <a href=\"https://github.com/HistoryLabs/events-api/blob/main/README.md\">here</a>."))
}
