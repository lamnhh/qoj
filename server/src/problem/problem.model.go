package problem

import (
	"database/sql"
	"errors"
	"fmt"
	"qoj/server/config"
	"strconv"
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

func CreateProblem(problem Problem) (int, error) {
	rows, err := config.DB.Query(
		"INSERT INTO problems(code, name, tl, ml) VALUES ($1, $2, $3, $4) RETURNING id",
		problem.Code,
		problem.Name,
		problem.TimeLimit,
		problem.MemoryLimit,
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

func FetchAllProblems(username string) ([]Problem, error) {
	rows, err := config.DB.Query("SELECT * FROM get_problem_list($1)", username)
	if err != nil {
		return []Problem{}, err
	}

	problemList := make([]Problem, 0)
	for rows.Next() {
		var problem Problem
		if err := rows.Scan(
			&problem.Id,
			&problem.Code,
			&problem.Name,
			&problem.TimeLimit,
			&problem.MemoryLimit,
			&problem.MaxScore,
			&problem.TestCount,
		); err != nil {
			return []Problem{}, err
		}
		normaliseProblem(&problem)
		problemList = append(problemList, problem)
	}

	return problemList, nil
}

func FetchProblemById(problemId int, username string) (Problem, error) {
	var problem Problem
	err := config.DB.QueryRow("SELECT * FROM get_problem_by_id($1, $2)", problemId, username).
		Scan(
			&problem.Id,
			&problem.Code,
			&problem.Name,
			&problem.TimeLimit,
			&problem.MemoryLimit,
			&problem.MaxScore,
			&problem.TestCount,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			// rows.Next() means no rows was returned. In other words, no problem with such ID exists
			return Problem{}, errors.New(fmt.Sprintf("No problem with ID %d exists", problemId))
		} else {
			return Problem{}, err
		}
	}

	normaliseProblem(&problem)
	return problem, nil
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