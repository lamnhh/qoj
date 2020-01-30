package test

import (
	"database/sql"
	"fmt"
	"os"
)

func GetFilePreview(path string) (string, error) {
	f, err := os.Open(path)
	defer func() { _ = f.Close() }()
	if err != nil {
		return "", err
	}

	buf := make([]byte, PREVIEW_SIZE)
	n, err := f.Read(buf)
	return string(buf[:n]), nil
}

func generateSingleValue(problemId int, order int, inputPreview string, outputPreview string) (string, []interface{}) {
	pos := order * 4 + 1
	key := fmt.Sprintf("($%d, $%d, $%d, $%d)", pos, pos + 1, pos + 2, pos + 3)
	val := []interface{}{problemId, order, inputPreview, outputPreview}
	return key, val
}

func parseTest(rows *sql.Rows) (Test, error) {
	var test Test
	err := rows.Scan(&test.Id, &test.ProblemId, &test.Order, &test.InputPreview, &test.OutputPreview)
	if err != nil {
		return Test{}, err
	}
	return test, nil
}

