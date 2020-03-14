package problem

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/config"
	"strconv"
)

func fetchProblemsSetBy(setter string, search string) ([]Problem, error) {
	cmd := `SELECT 
       	id, RTRIM(code), RTRIM(name), tl, ml
	FROM
		problems
	WHERE
	    setter = $1 AND (code ILIKE '%' || $2 || '%' or name ILIKE '%' || $2 || '%') ORDER BY id DESC`
	rows, err := config.DB.Query(cmd, setter, search)
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

func fetchProblemAdmin(problemId int) (Problem, error) {
	problem := Problem{}
	err := config.DB.QueryRow("SELECT id, RTRIM(code), RTRIM(name), tl, ml FROM problems WHERE id = $1", problemId).
		Scan(&problem.Id, &problem.Code, &problem.Name, &problem.TimeLimit, &problem.MemoryLimit)
	return problem, err
}

func getProblemAdmin(ctx *gin.Context) {
	setter := ctx.GetString("username")
	search := ctx.Query("search")
	problemList, err := fetchProblemsSetBy(setter, search)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, problemList)
	}
}

func getProblemIdAdmin(ctx *gin.Context) {
	problemId, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	problem, err := fetchProblemAdmin(int(problemId))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, problem)
}