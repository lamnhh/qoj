package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/src/token"
)

func getUser(ctx *gin.Context) {
	username := ctx.GetString("username")
	user, err := FindUserByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		user.Password = ""
		ctx.JSON(http.StatusOK, user)
	}
}

func InitialiseUserRoutes(app *gin.Engine) {
	app.GET("/api/user", token.RequireAuth(), getUser)
}