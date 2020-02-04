package submission

import (
	"errors"
	"fmt"
	"github.com/udhos/equalfile"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"qoj/server/src/problem"
	"qoj/server/src/queue"
	"qoj/server/src/test"
	"strings"
)

var judges map[int]chan interface{}

const (
	VerdictAc  = "Accepted"
	VerdictWa  = "Wrong Answer"
	VerdictRe  = "Runtime Error"
	VerdictTle = "Time Limit Exceeded"
	VerdictMle = "Memory Limit Exceeded"
)

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

	output, err := exec.Command(timeoutPath,
		"-t", fmt.Sprintf("%f", prob.TimeLimit),
		"-m", fmt.Sprintf("%d", prob.MemoryLimit * 1024),
		cmd,
	).CombinedOutput()
	if err == nil {
		result := strings.Split(string(output), " ")

		answerPreview := ""
		score := float32(0.0)
		verdict := ""
		switch result[0] {
			case "FINISHED":
				cmp := equalfile.New(nil, equalfile.Options{})
				equal, _ := cmp.CompareFile(outPath, tmpOutPath)

				answerPreview, _ = test.GetFilePreview(outPath)
				if equal {
					verdict = VerdictAc
					score = 1.0
				} else {
					verdict = VerdictWa
				}
			case "TIMEOUT":
				verdict = VerdictTle
			case "MEM":
				verdict = VerdictMle
			case "SIGNAL":
				verdict = VerdictRe
		}
		_ = updateScore(submissionId, testList[testId].Id, score, verdict, answerPreview)
	}

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

// TIMEOUT CPU 0.51 MEM 18612 MAXMEM 18612 STALE 0 MAXMEM_RSS 2500