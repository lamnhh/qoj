package result

import (
	"github.com/udhos/equalfile"
	"qoj/server/config"
	"qoj/server/src/test"
	"strconv"
	"strings"
)

type Result struct {
	InputPreview  string  `json:"inputPreview"`
	OutputPreview string  `json:"outputPreview"`
	AnswerPreview string  `json:"answerPreview"`
	Score         float32 `json:"score"`
	Verdict       string  `json:"verdict"`
	ExecutionTime float32 `json:"executionTime"`
	MemoryUsed    int     `json:"memoryUsed"`
}

const (
	VerdictAc  = "Accepted"
	VerdictWa  = "Wrong Answer"
	VerdictRe  = "Runtime Error"
	VerdictTle = "Time Limit Exceeded"
	VerdictMle = "Memory Limit Exceeded"
)

func ParseResultFromString(msg string, outPath string, ansPath string) Result {
	ans := Result{}

	// 0       1   2    3   4     5      6
	// TIMEOUT CPU 0.51 MEM 18612 MAXMEM 18612 STALE 0 MAXMEM_RSS 2500
	tokens := strings.Split(msg, " ")

	exeTime, _ := strconv.ParseFloat(tokens[2], 32)
	memUsed, _ := strconv.ParseInt(tokens[4], 10, 32)

	ans.ExecutionTime = float32(exeTime)
	ans.MemoryUsed = int(memUsed)
	ans.Score = 0
	switch tokens[0] {
	case "FINISHED":
		cmp := equalfile.New(nil, equalfile.Options{})
		equal, _ := cmp.CompareFile(outPath, ansPath)

		ans.AnswerPreview, _ = test.GetFilePreview(ansPath)
		if equal {
			ans.Verdict = VerdictAc
			ans.Score = 1.0
		} else {
			ans.Verdict = VerdictWa
		}
	case "TIMEOUT":
		ans.Verdict = VerdictTle
	case "MEM":
		ans.Verdict = VerdictMle
	case "SIGNAL":
		ans.Verdict = VerdictRe
	}

	return ans
}

func UpdateResult(submissionId int, testId int, result Result) error {
	cmd := `
	INSERT INTO
		submission_results
	VALUES
		($1, $2, $3, $4, $5, $6, $7)`
	_, err := config.DB.Exec(
		cmd,
		submissionId,
		testId,
		result.Score,
		result.Verdict,
		result.AnswerPreview,
		result.ExecutionTime,
		result.MemoryUsed,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetResultsOfSubmission(submissionId int, username string) ([]Result, error) {
	rows, err := config.DB.Query("SELECT * FROM get_submission_result($1, $2)", submissionId, username)
	if err != nil {
		return nil, err
	}

	resultList := make([]Result, 0)
	for rows.Next() {
		var result Result
		err := rows.Scan(
			&result.InputPreview,
			&result.OutputPreview,
			&result.AnswerPreview,
			&result.Score,
			&result.Verdict,
			&result.ExecutionTime,
			&result.MemoryUsed,
		)
		if err == nil {
			resultList = append(resultList, result)
		}
	}

	return resultList, nil
}