package contest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/src/token"
	"strconv"
)

func getContest(ctx *gin.Context) {
	contestList, err := fetchAllContests(ctx.GetString("username"))
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
	username := ctx.GetString("username")

	contest, err := fetchContestById(contestId, username)
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

func postContestIdRegister(ctx *gin.Context) {
	username := ctx.GetString("username")
	contestId64, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID"})
		return
	}
	contestId := int(contestId64)

	if err := joinContest(contestId, username); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func getContestIdParticipant(ctx *gin.Context) {
	contestId64, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID"})
		return
	}
	contestId := int(contestId64)

	participantList, err := fetchParticipantList(contestId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, participantList)
	}
}

func getContestIdProblem(ctx *gin.Context) {
	contestId64, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID"})
		return
	}
	contestId := int(contestId64)
	username := ctx.GetString("username")

	problemList, err := fetchProblemList(contestId, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, problemList)
	}
}

func getContestIdScore(ctx *gin.Context) {
	contestId64, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contest ID"})
		return
	}
	contestId := int(contestId64)

	scoreList, err := fetchContestScore(contestId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, scoreList)
	}
}

func InitialiseRoutes(app *gin.RouterGroup) {
	app.GET("/contest", token.ParseAuth(), getContest)
	app.GET("/contest/:id", token.ParseAuth(), getContestId)
	app.POST("/contest/:id/register", token.RequireAuth(), postContestIdRegister)
	app.GET("/contest/:id/participant", getContestIdParticipant)
	app.GET("/contest/:id/problem", token.ParseAuth(), getContestIdProblem)
	app.GET("/contest/:id/score", getContestIdScore)
}

func InitialiseAdminRoutes(app *gin.RouterGroup) {
	app.POST("/contest", token.RequireAuth(), postContest)
}
