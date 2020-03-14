package problem

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"qoj/server/config"
	"strconv"
	"strings"
)

type Problem struct {
	Id          int     `json:"id"`
	Code        string  `json:"code" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	TimeLimit   float32 `json:"timeLimit" binding:"required"`
	MemoryLimit int     `json:"memoryLimit" binding:"required"`
	MaxScore    float32 `json:"maxScore"`
	TestCount   int     `json:"testCount"`
}

func CreateProblem(problem Problem, setter string) (int, error) {
	rows, err := config.DB.Query(
		"INSERT INTO problems(code, name, tl, ml, setter) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		problem.Code,
		problem.Name,
		problem.TimeLimit,
		problem.MemoryLimit,
		setter,
	)
	if err != nil {
		return 0, err
	}

	var id int
	for rows.Next() {
		_ = rows.Scan(&id)
	}
	return id, nil
}

func DeleteProblem(problemId int) error {
	rows, err := config.DB.Query("DELETE FROM problems WHERE id = $1 RETURNING *", problemId)
	if err != nil {
		return err
	}
	if !rows.Next() {
		// rows.Next() means no rows was returned. In other words, no problem with such ID exists
		return errors.New(fmt.Sprintf("No problem with ID %d exists", problemId))
	}
	return nil
}

func FetchAllProblems(username string, page int, size int) ([]Problem, int, error) 	{
	rows, err := config.DB.Query("SELECT * FROM get_problem_list($1, $2, $3)", username, page, size)
	if err != nil {
		return []Problem{}, 0, err
	}

	problemList := make([]Problem, 0)
	for rows.Next() {
		problem, err := parseProblemFromRows(rows)
		if err == nil {
			problemList = append(problemList, problem)
		}
	}

	problemCount := 0
	_ = config.DB.QueryRow("SELECT COUNT(*) FROM problems").Scan(&problemCount)

	return problemList, problemCount, nil
}

func FetchProblemById(problemId int, username string) (Problem, error) {
	problemList, err := FetchProblemByIds([]int{problemId}, username)
	if err != nil {
		return Problem{}, err
	}
	if len(problemList) == 0 {
		return Problem{}, errors.New(fmt.Sprintf("Problem #%d does not exist", problemId))
	}
	return problemList[0], nil
}

func FetchProblemByIds(problemIds []int, username string) ([]Problem, error) {
	rows, err := config.DB.Query("SELECT * FROM get_problem_by_ids($1, $2)", pq.Array(problemIds), username)
	if err != nil {
		return []Problem{}, err
	}

	problemList := make([]Problem, 0)
	for rows.Next() {
		problem, err := parseProblemFromRows(rows)
		if err == nil {
			problemList = append(problemList, problem)
		}
	}

	return problemList, nil
}

func UpdateProblemMetadata(problemId int, patch map[string]string) (Problem, error) {
	// Fetch old problem metadata
	problem, err := FetchProblemById(problemId, "")
	if err != nil {
		return problem, err
	}

	// Update according to `patch`
	if val, ok := patch["code"]; ok {
		problem.Code = val
	}
	if val, ok := patch["name"]; ok {
		problem.Name = val
	}
	if val, ok := patch["timeLimit"]; ok {
		tl, _ := strconv.ParseFloat(val, 32)
		problem.TimeLimit = float32(tl)
	}
	if val, ok := patch["memoryLimit"]; ok {
		ml, _ := strconv.ParseInt(val, 10, 16)
		problem.MemoryLimit = int(ml)
	}

	// Update corresponding row in database
	err = config.DB.
		QueryRow("UPDATE problems SET code = $1, name = $2, tl = $3, ml = $4 WHERE id = $3 RETURNING *",
			problem.Code,
			problem.Name,
			problem.TimeLimit,
			problem.MemoryLimit,
			problem.Id,
		).
		Scan(&problem.Id, &problem.Code, &problem.Name, &problem.TimeLimit, &problem.MemoryLimit)
	if err != nil {
		return Problem{}, err
	}

	normaliseProblem(&problem)
	return problem, nil
}

func FetchSolvedProblemsOfUser(username string) ([]map[string]interface{}, error) {
	cmd := `
	SELECT
		DISTINCT(problem_id),
		problems.code
	FROM
		problems
		JOIN submissions ON (problems.id = submissions.problem_id)
		LEFT JOIN submission_results ON (submissions.id = submission_results.submission_id)
	WHERE
		submissions.username = $1
	GROUP BY
		submissions.id,
		problems.id,
		problems.code
	HAVING
		MIN(score) = 1
	ORDER BY
		problems.code, problem_id;`

	rows, err := config.DB.Query(cmd, username)
	if err != nil {
		return nil, err
	}

	problemList := make([]map[string]interface{}, 0)
	for rows.Next() {
		var (
			id   string
			code string
		)
		_ = rows.Scan(&id, &code)
		problemList = append(problemList, map[string]interface{}{
			"id":   id,
			"code": strings.TrimSpace(code),
		})
	}
	return problemList, nil
}

func FetchPartiallySolvedProblemsOfUser(username string) ([]map[string]interface{}, error) {
	cmd := `
	SELECT 
		id,
		code
	FROM
		(SELECT
			problems.*,
			COALESCE(MIN(score), 0) as score
		FROM
			problems
			JOIN submissions ON (problems.id = submissions.problem_id)
			LEFT JOIN submission_results ON (submissions.id = submission_results.submission_id)
		WHERE
			submissions.username = $1
		GROUP BY
			submissions.id, problems.id) s
	GROUP BY
		id,
		code
	HAVING
		MAX(score) < 1
	ORDER BY
		code;`

	rows, err := config.DB.Query(cmd, username)
	if err != nil {
		return nil, err
	}

	problemList := make([]map[string]interface{}, 0)
	for rows.Next() {
		var (
			id   string
			code string
		)
		_ = rows.Scan(&id, &code)
		problemList = append(problemList, map[string]interface{}{
			"id":   id,
			"code": strings.TrimSpace(code),
		})
	}
	return problemList, nil
}