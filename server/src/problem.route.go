package src

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getProblem(ctx *gin.Context) {
	problemList, err := FetchAllProblems()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, problemList)
}

func postProblem(ctx *gin.Context) {
	var problem Problem
	problem.Code = ctx.PostForm("code")
	problem.Name = ctx.PostForm("name")


}

func InitialiseProblemRoutes(app *gin.Engine) {
	app.GET("/api/problem", getProblem)
	app.POST("/api/problem", postProblem)
}
