package submission

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	submissionCount++
	submissionId := submissionCount

	// Initialise judge channel for this particular submission
	judges[submissionId] = make(chan interface{})
	listenerList[submissionId] = &ListenerList{}
	go submissionHandler(submissionId)

	problemId, err := strconv.ParseInt(ctx.PostForm("problemId"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = judge(submissionId, int(problemId), file)
	fmt.Println(err)

	ctx.JSON(http.StatusOK, gin.H{
		"submissionId": submissionId,
	})
}

func InitialiseSubmissionRoutes(app *gin.Engine) {
	submissionCount = 0

	app.POST("/api/submission", postSubmission)
}
