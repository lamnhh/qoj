package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitialiseApp() *gin.Engine {
	app := gin.Default()

	app.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Hello there",
		})
	})

	return app
}
