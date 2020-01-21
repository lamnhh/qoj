package src

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid2 "github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
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

	// Extract the zip file, verify filename
	extractedPath := filepath.Join(".", "server", "tasks", uuid)
	_, err := Unzip(zipPath, extractedPath)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H {"error": err.Error()})
		return
	}

	// Read all test subdirectories, each of which contains a test file
	fmt.Println(filepath.Join(extractedPath, problem.Code))
	testList, err := ReadDir(filepath.Join(extractedPath, problem.Code))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Iterate over all test directories, check if every directory contains <code>.INP and <code.OUT
	for _, testId := range testList {
		testPath := filepath.Join(extractedPath, problem.Code, "Test" + testId)
		if !DoesFileExists(filepath.Join(testPath, problem.Code + ".inp")) && !DoesFileExists(filepath.Join(testPath, problem.Code + ".INP")) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Test%s does not contain input file", problem.Code),
			})
			return
		}
		if !DoesFileExists(filepath.Join(testPath, problem.Code + ".out")) && !DoesFileExists(filepath.Join(testPath, problem.Code + ".OUT")) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Test%s does not contain output file", problem.Code),
			})
			return
		}
	}

	// Add problem to database to get problem ID
	problem.Id, err = CreateProblem(problem)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// All validation is completed
	// Rename uuid to problemId
	oldPath := filepath.Join(".", "server", "tasks", uuid)
	newPath := filepath.Join(".", "server", "tasks", fmt.Sprintf("%d", problem.Id))
	if err := os.Rename(oldPath, newPath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, problem)
}

func InitialiseProblemRoutes(app *gin.Engine) {
	app.GET("/api/problem", getProblem)
	app.POST("/api/problem", postProblem)
}
