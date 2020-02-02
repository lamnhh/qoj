package user

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/src/problem"
	"qoj/server/src/token"
)

func getUser(ctx *gin.Context) {
	username := ctx.Param("username")
	user, err := FindUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("User `%s` does not exist", username),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		user.Password = ""
		ctx.JSON(http.StatusOK, user)
	}
}

func getUserSolved(ctx *gin.Context) {
	username := ctx.Param("username")
	_, err := FindUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("User `%s` does not exist", username),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	problemList, err := problem.FetchSolvedProblemsOfUser(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, problemList)
	}
}

func getUserPartial(ctx *gin.Context) {
	username := ctx.Param("username")
	_, err := FindUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("User `%s` does not exist", username),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	problemList, err := problem.FetchPartiallySolvedProblemsOfUser(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, problemList)
	}
}

func InitialiseUserRoutes(app *gin.Engine) {
	app.GET("/api/user/:username", token.RequireAuth(), getUser)
	app.GET("/api/user/:username/solved", getUserSolved)
	app.GET("/api/user/:username/partial", getUserPartial)
}