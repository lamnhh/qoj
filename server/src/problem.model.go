package src

import (
	"database/sql"
	"errors"
	"fmt"
	"qoj/server/config"
	"strings"
)

type Problem struct {
	Id        int    `json:"id"`
	Code      string `json:"code" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

func CreateProblem(problem Problem) (int, error) {
	rows, err := config.DB.Query(
		"INSERT INTO problems(code, name) VALUES ($1, $2) RETURNING id",
		problem.Code,
		problem.Name,
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

func normaliseProblem(problem *Problem) {
	problem.Code = strings.TrimSpace(problem.Code)
	problem.Name = strings.TrimSpace(problem.Name)
}

func FetchAllProblems() ([]Problem, error) {
	rows, err := config.DB.Query("SELECT * FROM problems")
	if err != nil {
		return []Problem{}, err
	}

	problemList := make([]Problem, 0)
	for rows.Next() {
		var problem Problem
		if err := rows.Scan(&problem.Id, &problem.Code, &problem.Name); err != nil {
			return []Problem{}, err
		}
		normaliseProblem(&problem)
		problemList = append(problemList, problem)
	}

	return problemList, nil
}

func FetchProblemById(problemId int) (Problem, error) {
	var problem Problem
	err := config.DB.QueryRow("SELECT * FROM problems WHERE id = $1", problemId).
		Scan(&problem.Id, &problem.Code, &problem.Name)
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

func FetchProblemsByCode(code string) ([]Problem, error) {
	rows, err := config.DB.Query("SELECT * FROM problems WHERE code = $1", code)
	if err != nil {
		return []Problem{}, err
	}

	problemList := make([]Problem, 0)
	for rows.Next() {
		var problem Problem
		if err := rows.Scan(&problem.Id, &problem.Code, &problem.Name); err != nil {
			return []Problem{}, err
		}
		normaliseProblem(&problem)
		problemList = append(problemList, problem)
	}

	return problemList, nil
}

func updateProblemMetadata(problemId int, patch map[string]string) (Problem, error) {
	// Fetch old problem metadata
	problem, err := FetchProblemById(problemId)
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

	// Update corresponding row in database
	err = config.DB.
		QueryRow("UPDATE problems SET code = $1, name = $2 WHERE id = $3 RETURNING *", problem.Code, problem.Name, problem.Id).
		Scan(&problem.Id, &problem.Code, &problem.Name)
	if err != nil {
		return Problem{}, err
	}

	normaliseProblem(&problem)
	return problem, nil
}