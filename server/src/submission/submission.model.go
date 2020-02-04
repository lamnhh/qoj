package submission

import (
	"database/sql"
	"fmt"
	"qoj/server/config"
	"strings"
	"time"
)

type CodeSubmission struct {
	ProblemId int    `json:"problemId" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

type Submission struct {
	Id            int       `json:"id"`
	Username      string    `json:"username"`
	ProblemId     int       `json:"problemId"`
	ProblemName   string    `json:"problemName"`
	CreatedAt     time.Time `json:"createdAt"`
	Status        string    `json:"status"`
	ExecutionTime float32   `json:"executionTime"`
	MemoryUsed    int       `json:"memoryUsed"`
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
		&submission.ExecutionTime,
		&submission.MemoryUsed,
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
			submissions.status,
			COALESCE(MAX(execution_time), 0) as execution_time,
			COALESCE(MAX(memory_used), 0) as memory_used
		FROM
			submissions
			JOIN problems ON (submissions.problem_id = problems.id)
			LEFT JOIN submission_results ON (submissions.id = submission_results.submission_id)
		WHERE
			submissions.id = $1
		GROUP BY
			submissions.id,
			submissions.username,
			submissions.problem_id,
			problems.name,
			submissions.created_at,
			submissions.status`,
		submissionId,
	)
	if err != nil {
		return Submission{}, err
	}
	rows.Next()
	return parseSubmissionFromRow(rows)
}

func CountSubmission(filters map[string]interface{}) (int, error) {
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
		whereClause = "WHERE " + strings.Join(keyList, " AND ")
	}

	sql := fmt.Sprintf(`
	SELECT
		COUNT(*)
	FROM
		submissions
		JOIN problems ON (submissions.problem_id = problems.id)
	%s`, whereClause)

	ans := 0
	err := config.DB.QueryRow(sql, valList...).Scan(&ans)
	return ans, err
}

func FetchSubmissionList(filters map[string]interface{}, page int, size int) ([]Submission, error) {
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
		whereClause = "WHERE " + strings.Join(keyList, " AND ")
	}

	sql := fmt.Sprintf(`
	SELECT
		submissions.id,
		submissions.username,
		submissions.problem_id,
		problems.name,
		submissions.created_at,
		submissions.status,
		COALESCE(MAX(execution_time), 0) as execution_time,
		COALESCE(MAX(memory_used), 0) as memory_used
	FROM
		submissions
		JOIN problems ON (submissions.problem_id = problems.id)
		LEFT JOIN submission_results ON (submissions.id = submission_results.submission_id)
	%s
	GROUP BY
		submissions.id,
		submissions.username,
		submissions.problem_id,
		problems.name,
		submissions.created_at,
		submissions.status
	ORDER BY
		created_at DESC
	OFFSET %d LIMIT %d`, whereClause, page * size, size)

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