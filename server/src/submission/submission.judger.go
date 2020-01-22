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
	"strings"
)

type TestResult struct {
	Id int
	Result string
}

var judges map[int]chan TestResult

func judge(submissionId int, problemId int, fileHeader *multipart.FileHeader) error {
	dirname, _ := os.Getwd()

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

	// Compile the file above
	compileOutput, err := exec.Command("g++", cppPath, "-o", fmt.Sprintf("%d", submissionId)).Output()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Compile output:", string(compileOutput))

	judges[submissionId] = make(chan TestResult)
	go func() {
		timeoutPath := filepath.Join(dirname, "timeout")
		tmpOutPath := filepath.Join(dirname, fmt.Sprintf("%d.out", submissionId))
		exePath := filepath.Join(dirname, fmt.Sprintf("%d", submissionId))
		for testId := 1; ; testId++ {
			path := filepath.Join(dirname, "server", "tasks", fmt.Sprintf("%d", problemId))

			inpPath := filepath.Join(path, fmt.Sprintf("%d.inp", testId))
			outPath := filepath.Join(path, fmt.Sprintf("%d.out", testId))

			if !problem.DoesFileExists(inpPath) {
				judges[submissionId] <- TestResult{
					Id:     -1,
					Result: "",
				}
				break
			}

			fmt.Printf("Judging test %d\n", testId)

			cmd := fmt.Sprintf("%s < %s > %s", exePath, inpPath, tmpOutPath)

			output, err := exec.Command(timeoutPath, "-t", "1", cmd).CombinedOutput()
			if err != nil {
				log.Fatalln(err)
			}

			result := strings.Split(string(output), " ")
			resultMsg := ""
			switch result[0] {
			case "FINISHED":
				cmp := equalfile.New(nil, equalfile.Options{})
				equal, _ := cmp.CompareFile(outPath, tmpOutPath)
				if equal {
					resultMsg = "Correct"
				} else {
					resultMsg = "Wrong Answer"
				}
			case "TIMEOUT":
				resultMsg = "Time Limit Exceeded"
			case "MEM":
				resultMsg = "Memory Limit Exceeded"
			case "SIGNAL":
				resultMsg = "Runtime Error"
			default:
				continue
			}

			judges[submissionId] <- TestResult{
				Id:     testId,
				Result: resultMsg,
			}
		}

		fmt.Println("Done")
		_ = os.Remove(cppPath)
		_ = os.Remove(exePath)
		_ = os.Remove(tmpOutPath)
	}()
	return nil
}

// TIMEOUT CPU 0.51 MEM 18612 MAXMEM 18612 STALE 0 MAXMEM_RSS 2500