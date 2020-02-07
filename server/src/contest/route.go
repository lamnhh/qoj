package contest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/src/token"
	"strconv"
)

func getContest(ctx *gin.Context) {
	contestList, err := fetchAllContests()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, contestList)
	}
}

func getContestId(ctx *gin.Context) {
	contestId64, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID"})
		return
	}
	contestId := int(contestId64)

	contest, err := fetchContestById(contestId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, contest)
	}
}

func postContest(ctx *gin.Context) {
	body := Contest{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contest, err := createContest(body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, contest)
	}
}

func InitialiseContestRoutes(app *gin.Engine) {
	app.GET("/api/contest", getContest)
	app.GET("/api/contest/:id", getContestId)
	app.POST("/api/contest", token.RequireAuth(), postContest)
}
