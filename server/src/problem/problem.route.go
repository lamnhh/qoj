package problem

import (
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"net/http"
	"path/filepath"
	"qoj/server/src/test"
	"strconv"
)

func getProblem(ctx *gin.Context) {
	code := ctx.Query("code")

	var problemList []Problem
	var err error
	if code == "" {
		problemList, err = FetchAllProblems()
	} else {
		problemList, err = FetchProblemsByCode(code)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, problemList)
}

func postProblem(ctx *gin.Context) {
	var problem Problem

	// Parse problem code
	// Problem code is required. If form does not contain `code`, return a 400
	problem.Code = ctx.PostForm("code")
	if problem.Code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Problem Code required"})
		return
	}

	// Parse problem name
	// If form does not contain `name`, set `name` to the same value as `code`
	problem.Name = ctx.PostForm("name")
	if problem.Name == "" {
		problem.Name = problem.Code
	}

	// Parse problem's time limit
	tl, err := strconv.ParseFloat(ctx.PostForm("timeLimit"), 32)
	if err != nil || tl < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time limit"})
		return
	}
	problem.TimeLimit = float32(tl)

	// Parse problem's memory limit
	ml, err := strconv.ParseInt(ctx.PostForm("memoryLimit"), 10, 16)
	if err != nil || ml < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid memory limit"})
		return
	}
	problem.MemoryLimit = int(ml)

	// Parse test ZIP file
	// `file` is required
	file, _ := ctx.FormFile("file")
	if file == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ZIP file required"})
		return
	}

	// Save uploaded file to ./server/tasks/<uuid>.zip
	uuid := uuid2.New().String()
	zipPath := filepath.Join(".", "server", "tasks", uuid + ".zip")
	if err := ctx.SaveUploadedFile(file, zipPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Validate ZIP file
	code, err := validateTestZip(uuid, problem.Code)
	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
	} else {
		// Add problem to database to get problem ID
		problem.Id, err = CreateProblem(problem)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			// Create entries in table `tests`
			inpList, outList := saveTestData(uuid, problem.Id, problem.Code)
			_, _ = test.CreateTests(problem.Id, inpList, outList)
			ctx.JSON(http.StatusOK, problem)
		}
	}
	clearTemporaryData(uuid)
}

func getProblemId(ctx *gin.Context) {
	problemId, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	problem, err := FetchProblemById(int(problemId))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, problem)
}

func deleteProblemId(ctx *gin.Context) {
	problemId, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	if err := DeleteProblem(int(problemId)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func patchProblemId(ctx *gin.Context) {
	problemId, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	var patch map[string]string
	if err := ctx.ShouldBindJSON(&patch); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	problem, err := UpdateProblemMetadata(int(problemId), patch)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, problem)
}

func putProblemIdTest(ctx *gin.Context) {
	problemId, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	problem, err := FetchProblemById(int(problemId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	replace := ctx.Param("replace")
	if replace != "1" && replace != "0" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid `replace` param"})
		return
	}

	// Parse test ZIP file
	// `file` is required
	file, _ := ctx.FormFile("file")
	if file == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ZIP file required"})
		return
	}

	// Save uploaded file to ./server/tasks/<uuid>.zip
	uuid := uuid2.New().String()
	zipPath := filepath.Join(".", "server", "tasks", uuid + ".zip")
	if err := ctx.SaveUploadedFile(file, zipPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Validate ZIP file
	code, err := validateTestZip(uuid, problem.Code)
	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
	} else {
		// Replace old tests if user requires
		if replace == "1" {
			_ = test.DeleteAllTests(problem.Id)
		}

		// Create entries in table `tests`
		inpList, outList := saveTestData(uuid, problem.Id, problem.Code)
		_, _ = test.CreateTests(problem.Id, inpList, outList)
		ctx.JSON(http.StatusOK, problem)
	}
	clearTemporaryData(uuid)
}

func InitialiseProblemRoutes(app *gin.Engine) {
	app.GET("/api/problem", getProblem)
	app.GET("/api/problem/:id", getProblemId)
	app.POST("/api/problem", postProblem)
	app.DELETE("/api/problem/:id", deleteProblemId)
	app.PATCH("/api/problem/:id", patchProblemId)
	app.PUT("/api/problem/:id/test", putProblemIdTest)
}
