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
	Score       int       `json:"score"`
}

func parseSubmissionFromRow(rows *sql.Rows) (Submission, error) {
	var submission Submission
	err := rows.Scan(
		&submission.Id,
		&submission.Username,
		&submission.ProblemId,
		&submission.ProblemName,
		&submission.CreatedAt,
		&submission.Score,
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
			submissions.score	
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

	sql :=
		`SELECT
		submissions.id,
		submissions.username,
		submissions.problem_id,
		problems.name,
		submissions.created_at,
		submissions.score	
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

func updateScore(submissionId int, score int) error {
	_, err := config.DB.Exec("UPDATE submissions SET score = score + $1 WHERE id = $2", score, submissionId)
	return err
}
