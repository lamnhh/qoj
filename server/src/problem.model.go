package src

import (
	"errors"
	"fmt"
	"qoj/server/config"
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

func FetchAllProblems() ([]Problem, error) {
	rows, err := config.DB.Query("SELECT * FROM problems")
	if err != nil {
		return []Problem{}, err
	}

	var problemList []Problem
	for rows.Next() {
		var problem Problem
		if err := rows.Scan(&problem.Id, &problem.Code, &problem.Name); err != nil {
			return []Problem{}, err
		}
		problemList = append(problemList, problem)
	}

	return problemList, nil
}

func FetchProblemById(problemId int) (Problem, error) {
	rows, err := config.DB.Query("SELECT * FROM problems WHERE id = $1", problemId)
	if err != nil {
		return Problem{}, err
	}

	var problem Problem
	for rows.Next() {
		_ = rows.Scan(&problem.Id, &problem.Code, &problem.Name)
	}
	return problem, nil
}

func FetchProblemsByCode(code string) ([]Problem, error) {
	rows, err := config.DB.Query("SELECT * FROM problems WHERE code = $1", code)
	if err != nil {
		return []Problem{}, err
	}

	var problemList []Problem
	for rows.Next() {
		var problem Problem
		if err := rows.Scan(&problem.Id, &problem.Code, &problem.Name); err != nil {
			return []Problem{}, err
		}
		problemList = append(problemList, problem)
	}

	return problemList, nil
}