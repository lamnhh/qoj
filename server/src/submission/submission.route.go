package submission

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/src/auth"
	problem2 "qoj/server/src/problem"
	"strconv"
)

var submissionCount int

func submissionHandler(submissionId int) {
	for {
		select {
			case _res := <- judges[submissionId]:
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

	problem, err := problem2.FetchProblemById(problemId)
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
		filters["problemId"] = int(problemId)
	}

	if val := ctx.Query("username"); val != "" {
		filters["username"] = val
	}

	submissionList, err := fetchSubmissionList(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, submissionList)
	}
}

func InitialiseSubmissionRoutes(app *gin.Engine) {
	submissionCount = 0

	app.GET("/api/submission", getSubmission)
	app.POST("/api/submission", auth.RequireAuth(), postSubmission)
}
