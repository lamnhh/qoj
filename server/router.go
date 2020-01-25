package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/src/auth"
	"qoj/server/src/problem"
	"qoj/server/src/submission"
)

func InitialiseApp() *gin.Engine {
	app := gin.Default()

	app.Static("/static", "./static")
	app.Static("/node_modules", "./node_modules")
	app.LoadHTMLGlob("./templates/*")

	app.GET("/", func(ctx *gin.Context) {
		problemList, err := problem.FetchAllProblems()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.HTML(http.StatusOK, "index.tpl", problemList)
	})

	// Routing
	auth.InitialiseAuthRoutes(app)
	problem.InitialiseProblemRoutes(app)
	submission.InitialiseSubmissionSocket(app)
	submission.InitialiseSubmissionRoutes(app)

	return app
}
