package problem

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/config"
)

func fetchProblemsSetBy(setter string) ([]Problem, error) {
	cmd := "SELECT id, RTRIM(code), RTRIM(name), tl, ml FROM problems WHERE setter = $1 ORDER BY id DESC"
	rows, err := config.DB.Query(cmd, setter)
	if err != nil {
		return []Problem{}, err
	}

	problemList := make([]Problem, 0)
	for rows.Next() {
		problem := Problem{}
		err = rows.Scan(&problem.Id, &problem.Code, &problem.Name, &problem.TimeLimit, &problem.MemoryLimit)
		if err == nil {
			problemList = append(problemList, problem)
		}
	}

	return problemList, nil
}

func getProblemAdmin(ctx *gin.Context) {
	setter := ctx.GetString("username")
	problemList, err := fetchProblemsSetBy(setter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, problemList)
	}
}