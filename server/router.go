package server

import (
	"github.com/gin-gonic/gin"
	"qoj/server/src"
)

func InitialiseApp() *gin.Engine {
	app := gin.Default()

	// Routing
	src.InitialiseAuthRoutes(app)

	return app
}
