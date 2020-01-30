package submission

import (
	"fmt"
	"github.com/udhos/equalfile"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"qoj/server/src/problem"
	"qoj/server/src/queue"
	"qoj/server/src/test"
	"strings"
)

var judges map[int]chan interface{}

func judgeFunc(done chan interface{}, metadata interface{}) {
	config := metadata.(map[string]interface{})

	testId := config["testId"].(int)
	prob:= config["problem"].(problem.Problem)
	submissionId := config["submissionId"].(int)
	testList := config["testList"].([]test.Test)

	if testId >= len(testList) {
		// Clean up
		_ = os.Remove(fmt.Sprintf("%d", submissionId))
		_ = os.Remove(fmt.Sprintf("%d.cpp", submissionId))
		_ = os.Remove(fmt.Sprintf("%d.out", submissionId))
		done <- map[string]interface{}{
			"submissionId": submissionId,
			"type": "finish",
			"error": nil,
			"message": "",
		}
		return
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
	if err != nil {
		done <- map[string]interface{}{
			"submissionId": submissionId,
			"type": "result",
			"error": err,
			"message": "Runtime Error",
		}
	} else {
		result := strings.Split(string(output), " ")
		resultMsg := ""

		switch result[0] {
			case "FINISHED":
				cmp := equalfile.New(nil, equalfile.Options{})
				equal, _ := cmp.CompareFile(outPath, tmpOutPath)
				if equal {
					resultMsg = "Correct"
					_ = updateScore(submissionId, 1)
				} else {
					resultMsg = "Wrong Answer"
				}
			case "TIMEOUT":
				resultMsg = "Time Limit Exceeded"
			case "MEM":
				resultMsg = "Memory Limit Exceeded"
			case "SIGNAL":
				resultMsg = "Runtime Error"
		}

		done <- map[string]interface{}{
			"submissionId": submissionId,
			"type": "result",
			"error": nil,
			"message": fmt.Sprintf("%d | %s", testId, resultMsg),
		}
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

	cppPath := fmt.Sprintf("%d.cpp", submissionId)
	compileOutput, err := exec.Command("g++", cppPath, "-o", fmt.Sprintf("%d", submissionId)).CombinedOutput()
	if err != nil {
		// Compile error
		done <- map[string]interface{}{
			"submissionId": submissionId,
			"type":    "compile-error",
			"error":   err,
			"message": string(compileOutput),
		}
	} else {
		// Successfully compiled
		done <- map[string]interface{}{
			"submissionId": submissionId,
			"type":  "compile",
			"error": nil,
			"message": "",
		}

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

func judge(submissionId int, problem problem.Problem, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	// Save uploaded code to "{submissionId}.cpp"
	cppPath := fmt.Sprintf("%d.cpp", submissionId)
	cppFile, err := os.Create(cppPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(cppFile, file); err != nil {
		return err
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