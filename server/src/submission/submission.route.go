package submission

import (
	"fmt"
	"net/http"
	"qoj/server/src/common"
	problem2 "qoj/server/src/problem"
	"qoj/server/src/token"
	"strconv"

	"github.com/gin-gonic/gin"
)

var submissionCount int

func submissionHandler(submissionId int) {
	for {
		select {
		case _res := <-judges[submissionId]:
			res := _res.(map[string]interface{})

			connList := listenerList[submissionId].GetSubscriptionList()
			for _, conn := range connList {
				_ = conn.WriteJSON(res)
			}
			if res["type"] == "compile-error" || res["type"] == "finish" {
				return
			}
		}
	}
}

func postSubmission(ctx *gin.Context) {
	// Request must be a form-data

	// Parse current username
	usernameInterface, _ := ctx.Get("username")
	username := usernameInterface.(string)

	// Parse problemId
	problemIdInt64, err := strconv.ParseInt(ctx.PostForm("problemId"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	problemId := int(problemIdInt64)

	// Parse solution file
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create submission entry in database
	submission, err := createSubmission(username, problemId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	submissionId := submission.Id

	problem, err := problem2.FetchProblemById(problemId, "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Initialise judge channel for this particular submission
	judges[submissionId] = make(chan interface{})
	listenerList[submissionId] = &ListenerList{}
	go submissionHandler(submissionId)

	_ = judge(submissionId, problem, file)
	ctx.JSON(http.StatusOK, gin.H{
		"submissionId": submissionId,
	})
}

func getSubmission(ctx *gin.Context) {
	filters := make(map[string]interface{})

	if val := ctx.Query("problemId"); val != "" {
		problemId, err := strconv.ParseInt(val, 10, 16)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
			return
		}
		filters["problem_id"] = int(problemId)
	}

	if val := ctx.Query("username"); val != "" {
		filters["username"] = val
	}

	page := common.ParseQueryInt(ctx, "page", 1) - 1
	size := common.ParseQueryInt(ctx, "size", 20)
	submissionList, err := FetchSubmissionList(filters, page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count, err := CountSubmission(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		ctx.Writer.Header().Set("x-count", fmt.Sprintf("%d", count))
		ctx.JSON(http.StatusOK, submissionList)
	}
}

func getSubmissionIdResult(ctx *gin.Context) {
	submissionId64, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}
	submissionId := int(submissionId64)

	resultList, err := getSubmissionResults(submissionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resultList)
}

func InitialiseSubmissionRoutes(app *gin.Engine) {
	submissionCount = 0

	app.GET("/api/submission", getSubmission)
	app.POST("/api/submission", token.RequireAuth(), postSubmission)
	app.GET("/api/submission/:id/result", getSubmissionIdResult)
}
