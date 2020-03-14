package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/src/auth"
	"qoj/server/src/contest"
	"qoj/server/src/language"
	"qoj/server/src/problem"
	"qoj/server/src/submission"
	"qoj/server/src/user"
)

func InitialiseApp() *gin.Engine {
	app := gin.Default()

	app.Static("/static", "./static")
	app.Static("/node_modules", "./node_modules")
	app.Static("/profile-picture", "./server/profile-picture")
	app.LoadHTMLGlob("./static/*.html")

	// Routing
	client := app.Group("/api")
	{
		auth.InitialiseRoutes(client)
		problem.InitialiseRoutes(client)
		submission.InitialiseRoutes(client)
		user.InitialiseRoutes(client)
		language.InitialiseRoutes(client)
		contest.InitialiseRoutes(client)
	}
	submission.InitialiseSocket(app)
	contest.InitialiseSocket(app)

	// Serve client js
	app.Use(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", []string{
			"/node_modules/codemirror/lib/codemirror.js",
			"/node_modules/codemirror/mode/clike/clike.js",
			"/static/index.js",
		})
	})

	return app
}

func InitialiseAdminApp() *gin.Engine {
	app := gin.Default()

	app.Static("/static", "./static")
	app.Static("/node_modules", "./node_modules")
	app.LoadHTMLGlob("./static/*.html")

	// Admin routing
	admin := app.Group("/api")
	{
		auth.InitialiseAdminRoutes(admin)
		user.InitialiseAdminRoutes(admin)
		problem.InitialiseAdminRoutes(admin)
		contest.InitialiseAdminRoutes(admin)
	}

	// Serve admin js
	app.Use(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", []string{
			"/static/admin.js",
		})
	})

	return app
}
