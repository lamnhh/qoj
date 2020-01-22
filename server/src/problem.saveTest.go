package src

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// validateTestZip validates the uploaded ZIP file
// Parameters:
// - uuid: a unique ID for the uploaded file
// - problemCode: code (i.e. NKPALIN) of the to-be-uploaded problem
// This function assumes that the ZIP file has been save to "./server/tasks/uuid.zip"
func validateTestZip(uuid string, problemCode string) (int, error) {
	// Extract the zip file, verify filename
	zipPath := filepath.Join(".", "server", "tasks", uuid + ".zip")
	extractedPath := filepath.Join(".", "server", "tasks", uuid)
	_, err := Unzip(zipPath, extractedPath)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// Read all test subdirectories, each of which contains a test file
	testList, err := ReadDir(filepath.Join(extractedPath, problemCode))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Iterate over all test directories, check if every directory contains <code>.INP and <code.OUT
	for _, testId := range testList {
		// fileName = ".../Code/Test00/Code.(INP|inp|OUT|out)
		fileName := filepath.Join(extractedPath, problemCode, "Test" + testId, problemCode)

		// Check if "Code.INP" or "Code.inp" exists
		if !DoesFileExists(fileName + ".INP") && !DoesFileExists(fileName + ".inp") {
			message := fmt.Sprintf("Test%s does not contain input file", problemCode)
			return http.StatusBadRequest, errors.New(message)
		}

		// Check if "Code.OUT" or "Code.out" exists
		if !DoesFileExists(fileName + ".OUT") && !DoesFileExists(fileName + ".out") {
			message := fmt.Sprintf("Test%s does not contain output file", problemCode)
			return http.StatusBadRequest, errors.New(message)
		}
	}

	return 0, nil
}

func saveTestData(uuid string, problemId int, problemCode string) {
	extractedPath := filepath.Join(".", "server", "tasks", uuid)

	// Path to save test data: ./server/tasks/id/*.inp | *.out
	testPath := filepath.Join(".", "server", "tasks", fmt.Sprintf("%d", problemId))
	_ = os.RemoveAll(testPath)
	_ = os.Mkdir(testPath, os.ModePerm)

	// Read all test subdirectories, each of which contains a test file
	testList, _ := ReadDir(filepath.Join(extractedPath, problemCode))

	// Iterate over all test directories, check if every directory contains <code>.INP and <code.OUT
	for _, testId := range testList {
		// fileName = ".../Code/Test00/Code.(INP|inp|OUT|out)
		fileName := filepath.Join(extractedPath, problemCode, "Test" + testId, problemCode)

		// Fetch inputFile path. This file should be in the form "fileName.INP" or "fileName.inp"
		var inputFile string
		if DoesFileExists(fileName + ".INP") {
			inputFile = fileName + ".INP"
		} else {
			inputFile = fileName + ".inp"
		}

		// Fetch outputFile path. This file should be in the form "fileName.OUT" or "fileName.out"
		var outputFile string
		if DoesFileExists(fileName + ".OUT") {
			outputFile = fileName + ".OUT"
		} else {
			outputFile = fileName + ".out"
		}

		// Save both files to "testId.inp", "testId.out"
		_ = os.Rename(inputFile, filepath.Join(testPath, testId + ".inp"))
		_ = os.Rename(outputFile, filepath.Join(testPath, testId + ".out"))
	}
}

func clearTemporaryData(uuid string) {
	_ = os.RemoveAll(filepath.Join(".", "server", "tasks", uuid))
	_ = os.Remove(filepath.Join(".", "server", "tasks", uuid + ".zip"))
}