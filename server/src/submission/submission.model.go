package submission

import (
	"database/sql"
	"fmt"
	"qoj/server/config"
	"strings"
	"time"
)

type Submission struct {
	Id          int       `json:"id"`
	Username    string    `json:"username"`
	ProblemId   int       `json:"problemId"`
	ProblemName string    `json:"problemName"`
	CreatedAt   time.Time `json:"createdAt"`
	Status      string    `json:"status"`
}

type SubmissionResult struct {
	InputPreview  string  `json:"inputPreview"`
	OutputPreview string  `json:"outputPreview"`
	AnswerPreview string  `json:"answerPreview"`
	Score         float32 `json:"score"`
	Verdict       string  `json:"verdict"`
}

func parseSubmissionFromRow(rows *sql.Rows) (Submission, error) {
	var submission Submission
	err := rows.Scan(
		&submission.Id,
		&submission.Username,
		&submission.ProblemId,
		&submission.ProblemName,
		&submission.CreatedAt,
		&submission.Status,
	)
	if err != nil {
		return Submission{}, err
	}

	submission.Username = strings.TrimSpace(submission.Username)
	submission.ProblemName = strings.TrimSpace(submission.ProblemName)
	return submission, nil
}

func createSubmission(username string, problemId int) (Submission, error) {
	rows, err := config.DB.Query(
		"INSERT INTO submissions(username, problem_id) VALUES ($1, $2) RETURNING id",
		username,
		problemId,
	)
	if err != nil {
		return Submission{}, err
	}

	var submissionId int
	for rows.Next() {
		_ = rows.Scan(&submissionId)
	}
	return fetchSubmissionById(submissionId)
}

func fetchSubmissionById(submissionId int) (Submission, error) {
	rows, err := config.DB.Query(
		`SELECT
			submissions.id,
			submissions.username,
			submissions.problem_id,
			problems.name,
			submissions.created_at,
			submissions.status	
		FROM
			submissions
			JOIN problems ON (submissions.problem_id = problems.id)
		WHERE
			submissions.id = $1`,
		submissionId,
	)
	if err != nil {
		return Submission{}, err
	}
	rows.Next()
	return parseSubmissionFromRow(rows)
}

func FetchSubmissionList(filters map[string]interface{}) ([]Submission, error) {
	keyList := make([]string, 0)
	valList := make([]interface{}, 0)
	count := 0
	for k, v := range filters {
		count++
		keyList = append(keyList, fmt.Sprintf("%s = $%d", k, count))
		valList = append(valList, v)
	}

	var whereClause string
	if len(keyList) == 0 {
		whereClause = ""
	} else {
		whereClause = "WHERE " + strings.Join(keyList, ", ")
	}

	sql := `
	SELECT
		submissions.id,
		submissions.username,
		submissions.problem_id,
		problems.name,
		submissions.created_at,
		submissions.status	
	FROM
		submissions
		JOIN problems ON (submissions.problem_id = problems.id)` +
		whereClause +
	`ORDER BY
		created_at DESC`

	rows, err := config.DB.Query(sql, valList...)
	if err != nil {
		return []Submission{}, err
	}

	submissionList := make([]Submission, 0)
	for rows.Next() {
		submission, err := parseSubmissionFromRow(rows)
		if err != nil {
			return []Submission{}, err
		}
		submissionList = append(submissionList, submission)
	}
	return submissionList, nil
}

func updateSubmissionStatus(submissionId int, status string) error {
	_, err := config.DB.Exec("UPDATE submissions SET status = $1 WHERE id = $2", status, submissionId)
	return err
}

func getSubmissionScore(submissionId int) float32 {
	score := float32(0)
	err := config.DB.
		QueryRow("SELECT SUM(score) FROM submission_results WHERE submission_id = $1", submissionId).
		Scan(&score)
	if err != nil {
		score = 0
	}
	return score
}

func getSubmissionResults(submissionId int) ([]SubmissionResult, error) {
	rows, err := config.DB.Query("SELECT * FROM get_submission_result($1)", submissionId)
	if err != nil {
		return nil, err
	}

	resultList := make([]SubmissionResult, 0)
	for rows.Next() {
		var result SubmissionResult
		err := rows.Scan(
			&result.InputPreview,
			&result.OutputPreview,
			&result.AnswerPreview,
			&result.Score,
			&result.Verdict,
		)
		if err == nil {
			resultList = append(resultList, result)
		}
	}

	return resultList, nil
}