package submission

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"qoj/server/src/problem"
	"qoj/server/src/queue"
	result2 "qoj/server/src/result"
	"qoj/server/src/test"
)

var judges map[int]chan interface{}

func judgeFunc(done chan interface{}, metadata interface{}) {
	config := metadata.(map[string]interface{})

	testId := config["testId"].(int)
	prob:= config["problem"].(problem.Problem)
	submissionId := config["submissionId"].(int)
	testList := config["testList"].([]test.Test)

	if testId >= len(testList) {
		status := fmt.Sprintf("%.2f / %d.00", getSubmissionScore(submissionId), len(testList))
		_ = updateSubmissionStatus(submissionId, status)
		// Clean up
		_ = os.Remove(fmt.Sprintf("%d", submissionId))
		_ = os.Remove(fmt.Sprintf("%d.cpp", submissionId))
		_ = os.Remove(fmt.Sprintf("%d.out", submissionId))
		done <- map[string]interface{}{
			"submissionId": submissionId,
			"message":      status,
		}
		return
	}

	status := fmt.Sprintf("Running on test %d...", testId + 1)
	_ = updateSubmissionStatus(submissionId, status)
	done <- map[string]interface{}{
		"submissionId": submissionId,
		"message": status,
	}

	dirname, _ := os.Getwd()
	timeoutPath := filepath.Join(dirname, "timeout")
	tmpOutPath := filepath.Join(dirname, fmt.Sprintf("%d.out", submissionId))
	exePath := filepath.Join(dirname, fmt.Sprintf("%d", submissionId))

	inpPath := filepath.Join(dirname, "server", "tasks", fmt.Sprintf("%d.inp", testList[testId].Id))
	outPath := filepath.Join(dirname, "server", "tasks", fmt.Sprintf("%d.out", testList[testId].Id))

	log.Printf("Judging test %d\n", testId)

	cmd := fmt.Sprintf("%s < %s > %s", exePath, inpPath, tmpOutPath)

	output, _ := exec.Command(timeoutPath,
		"-t", fmt.Sprintf("%f", prob.TimeLimit),
		"-m", fmt.Sprintf("%d", prob.MemoryLimit * 1024),
		cmd,
	).CombinedOutput()

	result := result2.ParseResultFromString(string(output), tmpOutPath, outPath)
	_ = result2.UpdateResult(submissionId, testList[testId].Id, result)

	queue.Push(queue.Task{
		Run:           judgeFunc,
		ResultChannel: done,
		Params: map[string]interface{}{
			"testId":       testId + 1,
			"problem":      prob,
			"submissionId": submissionId,
			"testList":     testList,
		},
	})
}

func compileFunc(done chan interface{}, metadata interface{}) {
	config := metadata.(map[string]interface{})

	prob := config["problem"].(problem.Problem)
	submissionId := config["submissionId"].(int)

	// Update status
	done <- map[string]interface{}{
		"submissionId": submissionId,
		"message":      "Compiling...",
	}
	_ = updateSubmissionStatus(submissionId, "Compiling...")

	cppPath := fmt.Sprintf("%d.cpp", submissionId)
	compileOutput, err := exec.Command("g++", cppPath, "-o", fmt.Sprintf("%d", submissionId)).CombinedOutput()
	if err != nil {
		// Compile error
		done <- map[string]interface{}{
			"submissionId": submissionId,
			"message":      "Compile Error|" + string(compileOutput),
		}
		_ = updateSubmissionStatus(submissionId, "Compile Error|" + string(compileOutput))
		_ = os.Remove(cppPath)
	} else {
		// Successfully compiled
		testList, _ := test.FetchAllTests(prob.Id)
		queue.Push(queue.Task{
			Run:           judgeFunc,
			ResultChannel: done,
			Params: map[string]interface{}{
				"testId":       0,
				"problem":      prob,
				"submissionId": submissionId,
				"testList":     testList,
			},
		})
	}
}

func judge(submissionId int, problem problem.Problem, code string) error {
	// Create file "{submissionId}.cpp"
	file, err := os.Create(fmt.Sprintf("%d.cpp", submissionId))
	if err != nil {
		return err
	}

	// Write `code` to `file`
	n, err := file.WriteString(code)
	if err != nil {
		return err
	}
	if n != len(code) {
		return errors.New("Internal Server Error")
	}

	compileTask := queue.Task{
		Run:           compileFunc,
		ResultChannel: judges[submissionId],
		Params:       map[string]interface{}{
			"problem":      problem,
			"submissionId": submissionId,
		},
	}

	queue.Push(compileTask)
	return nil
}