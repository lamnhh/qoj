package server

import (
	"github.com/gin-gonic/gin"
	"qoj/server/src/auth"
	"qoj/server/src/problem"
)

func InitialiseApp() *gin.Engine {
	app := gin.Default()

	// Routing
	auth.InitialiseAuthRoutes(app)
	problem.InitialiseProblemRoutes(app)

	return app
}
