package test

import (
	"fmt"
	"os"
	"path/filepath"
	"qoj/server/config"
	"strings"
)

const PREVIEW_SIZE int = 100

type Test struct {
	Id            int    `json:"id"`
	ProblemId     int    `json:"problemId"`
	Order         int    `json:"order"`
	InputPreview  string `json:"inputPreview"`
	OutputPreview string `json:"outputPreview"`
}

func CreateTests(problemId int, tmpInputPath []string, tmpOutputPath []string) ([]Test, error) {
	keyList := make([]string, 0)
	valList := make([]interface{}, 0)

	testCount := len(tmpInputPath)
	for i := 0; i < testCount; i++ {
		inputPreview, _ := GetFilePreview(tmpInputPath[i])
		outputPreview, _ := GetFilePreview(tmpOutputPath[i])

		key, val := generateSingleValue(problemId, i, inputPreview, outputPreview)
		keyList = append(keyList, key)
		valList = append(valList, val...)
	}

	sql := "INSERT INTO tests(problem_id, ord, inp_preview, out_preview) VALUES " + strings.Join(keyList, ", ") + " RETURNING *"
	rows, err := config.DB.Query(sql, valList...)
	if err != nil {
		return []Test{}, err
	}

	testList := make([]Test, 0)
	for rows.Next() {
		test, _ := parseTest(rows)
		testList = append(testList, test)
	}

	targetTestPath := filepath.Join(".", "server", "tasks")
	for i, test := range testList {
		_ = os.Rename(tmpInputPath[i], filepath.Join(targetTestPath, fmt.Sprintf("%d.inp", test.Id)))
		_ = os.Rename(tmpOutputPath[i], filepath.Join(targetTestPath, fmt.Sprintf("%d.out", test.Id)))
	}

	return testList, nil
}

func FetchAllTests(problemId int) ([]Test, error) {
	// &test.Id, &test.ProblemId, &test.Order, &test.InputPreview, &test.OutputPreview
	cmd := `
	SELECT
		tests.id, tests.problem_id, tests.ord, tests.inp_preview, tests.out_preview
	FROM
		tests
		JOIN problems ON (tests.problem_id = problems.original_id)
	WHERE
		problems.id = $1
	ORDER BY 
		ord ASC`
	rows, err := config.DB.Query(cmd, problemId)
	if err != nil {
		return []Test{}, err
	}

	testList := make([]Test, 0)
	for rows.Next() {
		test, _ := parseTest(rows)
		testList = append(testList, test)
	}
	return testList, nil
}

func DeleteAllTests(problemId int) error {
	_, err := config.DB.Exec("DELETE FROM tests WHERE problem_id = $1", problemId)
	return err
}