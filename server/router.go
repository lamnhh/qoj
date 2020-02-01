package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/src/auth"
	"qoj/server/src/problem"
	"qoj/server/src/submission"
	"qoj/server/src/user"
)

func InitialiseApp() *gin.Engine {
	app := gin.Default()

	app.Static("/static", "./static")
	app.Static("/node_modules", "./node_modules")
	app.LoadHTMLGlob("./static/*.html")

	// Routing
	auth.InitialiseAuthRoutes(app)
	problem.InitialiseProblemRoutes(app)
	submission.InitialiseSubmissionSocket(app)
	submission.InitialiseSubmissionRoutes(app)
	user.InitialiseUserRoutes(app)

	app.Use(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	return app
}
