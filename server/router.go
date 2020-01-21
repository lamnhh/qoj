package server

import (
	"github.com/gin-gonic/gin"
	"qoj/server/src"
)

func InitialiseApp() *gin.Engine {
	app := gin.Default()

	// Routing
	src.InitialiseAuthRoute(app)

	return app
}
