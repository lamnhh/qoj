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
	"strings"
)

var judges map[int]chan interface{}

func judgeFunc(done chan interface{}, metadata interface{}) {
	config := metadata.(map[string]interface{})

	testId := config["testId"].(int)
	prob:= config["problem"].(problem.Problem)
	submissionId := config["submissionId"].(int)

	dirname, _ := os.Getwd()
	timeoutPath := filepath.Join(dirname, "timeout")
	tmpOutPath := filepath.Join(dirname, fmt.Sprintf("%d.out", submissionId))
	exePath := filepath.Join(dirname, fmt.Sprintf("%d", submissionId))

	path := filepath.Join(dirname, "server", "tasks", fmt.Sprintf("%d", prob.Id))

	inpPath := filepath.Join(path, fmt.Sprintf("%d.inp", testId))
	outPath := filepath.Join(path, fmt.Sprintf("%d.out", testId))

	if !problem.DoesFileExists(inpPath) {
		// Clean up
		_ = os.Remove(fmt.Sprintf("%d", submissionId))
		_ = os.Remove(fmt.Sprintf("%d.cpp", submissionId))
		_ = os.Remove(fmt.Sprintf("%d.out", submissionId))
		done <- map[string]interface{}{
			"type": "finish",
			"error": nil,
			"message": "",
		}
		return
	}

	log.Printf("Judging test %d\n", testId)

	cmd := fmt.Sprintf("%s < %s > %s", exePath, inpPath, tmpOutPath)

	output, err := exec.Command(timeoutPath,
		"-t", fmt.Sprintf("%f", prob.TimeLimit),
		"-m", fmt.Sprintf("%d", prob.MemoryLimit * 1024),
		cmd,
	).CombinedOutput()
	if err != nil {
		done <- map[string]interface{}{
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
			"type":    "compile-error",
			"error":   err,
			"message": string(compileOutput),
		}
	} else {
		// Successfully compiled
		done <- map[string]interface{}{
			"type":  "compile",
			"error": nil,
		}
		queue.Push(queue.Task{
			Run:           judgeFunc,
			ResultChannel: done,
			Params: map[string]interface{}{
				"testId":       1,
				"problem":      prob,
				"submissionId": submissionId,
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