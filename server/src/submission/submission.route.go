package submission

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"qoj/server/src/common"
	"qoj/server/src/language"
	problem2 "qoj/server/src/problem"
	"qoj/server/src/result"
	"qoj/server/src/token"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	username := ctx.GetString("username")
	body := CodeSubmission{}

	// Parse JSON body
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate code length (<= 50000B)
	if len(body.Code) > 50000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Solution file exceeds 50000B"})
		return
	}

	// Check if problemId exists
	problem, err := problem2.FetchProblemById(body.ProblemId, "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if languageId exists
	lang, err := language.FetchLanguageById(body.LanguageId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create submission entry in database
	submission, err := createSubmission(username, body.ProblemId, body.LanguageId, body.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	submissionId := submission.Id

	// Initialise judge channel for this particular submission
	judges[submissionId] = make(chan interface{})
	listenerList[submissionId] = &ListenerList{}
	go submissionHandler(submissionId)

	_ = judge(submissionId, body.Code, problem, lang)
	ctx.JSON(http.StatusOK, gin.H{
		"submissionId": submissionId,
	})
}

func getSubmission(ctx *gin.Context) {
	filters := make(map[string]interface{})

	if val := ctx.QueryArray("problemId"); len(val) != 0 {
		problemIds := make([]int, 0)
		for _, id := range val {
			problemId, err := strconv.ParseInt(id, 10, 16)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
				return
			}
			problemIds = append(problemIds, int(problemId))
		}
		filters["problem_id"] = pq.Array(problemIds)
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

func getSubmissionId(ctx *gin.Context) {
	submissionId64, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}
	submissionId := int(submissionId64)

	submission, err := FetchSubmissionById(submissionId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Submission #%d does not exist", submissionId)})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		ctx.JSON(http.StatusOK, submission)
	}
}

func getSubmissionIdResult(ctx *gin.Context) {
	submissionId64, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}
	submissionId := int(submissionId64)

	resultList, err := result.GetResultsOfSubmission(submissionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resultList)
}

func getSubmissionIdCode(ctx *gin.Context) {
	submissionId64, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}
	submissionId := int(submissionId64)

	code, err := fetchCode(submissionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"code": code})
	}
}

func getSubmissionIdCompile(ctx *gin.Context) {
	submissionId64, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}
	submissionId := int(submissionId64)

	msg, err := fetchCompilationMessage(submissionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"compileMessage": msg})
	}
}

func InitialiseSubmissionRoutes(app *gin.Engine) {
	app.GET("/api/submission", getSubmission)
	app.GET("/api/submission/:id", getSubmissionId)
	app.POST("/api/submission", token.RequireAuth(), postSubmission)
	app.GET("/api/submission/:id/result", getSubmissionIdResult)
	app.GET("/api/submission/:id/code", getSubmissionIdCode)
	app.GET("/api/submission/:id/compile", getSubmissionIdCompile)
}