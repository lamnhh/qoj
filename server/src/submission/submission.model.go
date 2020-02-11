package submission

import (
	"database/sql"
	"fmt"
	"qoj/server/config"
	"strings"
	"time"
)

type CodeSubmission struct {
	ProblemId  int    `json:"problemId" binding:"required"`
	Code       string `json:"code" binding:"required"`
	LanguageId int    `json:"languageId" binding:"required"`
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
	LanguageId    int       `json:"languageId"`
	Language      string    `json:"language"`
	Code          string    `json:"code,omitempty"`
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
		&submission.LanguageId,
		&submission.Language,
	)
	if err != nil {
		return Submission{}, err
	}

	submission.Username = strings.TrimSpace(submission.Username)
	submission.ProblemName = strings.TrimSpace(submission.ProblemName)
	return submission, nil
}

func createSubmission(username string, problemId int, languageId int, code string) (Submission, error) {
	rows, err := config.DB.Query(
		"INSERT INTO submissions(username, problem_id, language_id, code) VALUES ($1, $2, $3, $4) RETURNING id",
		username,
		problemId,
		languageId,
		code,
	)
	if err != nil {
		return Submission{}, err
	}

	var submissionId int
	for rows.Next() {
		_ = rows.Scan(&submissionId)
	}
	return FetchSubmissionById(submissionId)
}

func updateCompilationMessage(submissionId int, msg string) error {
	_, err := config.DB.Exec("UPDATE submissions SET compile_msg = $1 WHERE id = $2", msg, submissionId)
	return err
}

func fetchCode(submissionId int) (string, error) {
	code := ""
	err := config.DB.QueryRow("SELECT code FROM submissions WHERE id = $1", submissionId).Scan(&code)
	return code, err
}

func fetchCompilationMessage(submissionId int) (string, error) {
	msg := ""
	err := config.DB.QueryRow("SELECT compile_msg FROM submissions WHERE id = $1", submissionId).Scan(&msg)
	return msg, err
}

func FetchSubmissionById(submissionId int) (Submission, error) {
	rows, err := config.DB.Query(
		`SELECT
			submissions.id,
			submissions.username,
			submissions.problem_id,
			problems.name,
			submissions.created_at,
			submissions.status,
			COALESCE(MAX(execution_time), 0) as execution_time,
			COALESCE(MAX(memory_used), 0) as memory_used,
			languages.id,
			languages.name
		FROM
			submissions
			JOIN problems ON (submissions.problem_id = problems.id)
			JOIN languages ON (submissions.language_id = languages.id)
			LEFT JOIN submission_results ON (submissions.id = submission_results.submission_id)
		WHERE
			submissions.id = $1
		GROUP BY
			submissions.id,
			submissions.username,
			submissions.problem_id,
			problems.name,
			submissions.created_at,
			submissions.status,
			languages.id,
			languages.name`,
		submissionId,
	)
	if err != nil {
		return Submission{}, err
	}

	var ans Submission
	for rows.Next() {
		ans, _ = parseSubmissionFromRow(rows)
	}
	return ans, nil
}

func CountSubmission(filters map[string]interface{}, allowInContest bool) (int, error) {
	keyList := make([]string, 0)
	valList := make([]interface{}, 0)
	count := 0
	for k, v := range filters {
		count++
		if k != "problem_id" {
			keyList = append(keyList, fmt.Sprintf("%s = ($%d)", k, count))
		} else {
			keyList = append(keyList, fmt.Sprintf("%s = ANY($%d)", k, count))
		}
		valList = append(valList, v)
	}

	if allowInContest == false {
		keyList = append(keyList, "problems.contest_id IS NULL")
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

func FetchSubmissionList(filters map[string]interface{}, allowInContest bool, page int, size int) ([]Submission, error) {
	keyList := make([]string, 0)
	valList := make([]interface{}, 0)
	count := 0
	for k, v := range filters {
		count++
		if k != "problem_id" {
			keyList = append(keyList, fmt.Sprintf("%s = ($%d)", k, count))
		} else {
			keyList = append(keyList, fmt.Sprintf("%s = ANY($%d)", k, count))
		}
		valList = append(valList, v)
	}

	if allowInContest == false {
		keyList = append(keyList, "problems.contest_id IS NULL")
	}

	var whereClause string
	if len(keyList) == 0 {
		whereClause = ""
	} else {
		whereClause = "WHERE " + strings.Join(keyList, " AND ")
	}

	cmd := fmt.Sprintf(`
	SELECT
		submissions.id,
		submissions.username,
		submissions.problem_id,
		problems.name,
		submissions.created_at,
		submissions.status,
		COALESCE(MAX(execution_time), 0) as execution_time,
		COALESCE(MAX(memory_used), 0) as memory_used,
		languages.id,
		RTRIM(languages.name)
	FROM
		submissions
		JOIN problems ON (submissions.problem_id = problems.id)
		JOIN languages ON (submissions.language_id = languages.id)
		LEFT JOIN submission_results ON (submissions.id = submission_results.submission_id)
	%s
	GROUP BY
		submissions.id,
		submissions.username,
		submissions.problem_id,
		problems.name,
		submissions.created_at,
		submissions.status,
		languages.id,
		languages.name
	ORDER BY
		created_at DESC`, whereClause)

	if size != -1 {
		cmd = cmd + fmt.Sprintf(` OFFSET %d LIMIT %d`, page*size, size)
	}

	rows, err := config.DB.Query(cmd, valList...)
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