package contest

import (
	"database/sql"
	"github.com/lib/pq"
	"qoj/server/src/problem"
)

func int64toInt(a []int64) []int {
	b := make([]int, 0)
	for _, x := range a {
		b = append(b, int(x))
	}
	return b
}

func parseSingleContest(row *sql.Row, username string) (SpecifiedContest, error) {
	contest := SpecifiedContest{}
	ids64 := make([]int64, 0)
	err := row.Scan(&contest.Id, &contest.Name, pq.Array(&ids64), &contest.StartDate, &contest.Duration)
	if err != nil {
		return contest, err
	}

	ids := make([]int, 0)
	for _, x := range ids64 {
		ids = append(ids, int(x))
	}
	contest.ProblemList, err = problem.FetchProblemByIds(ids, username)
	return contest, err
}

func parseMultipleContests(rows *sql.Rows) (Contest, error) {
	contest := Contest{}
	ids := make([]int64, 0)
	err := rows.Scan(&contest.Id, &contest.Name, pq.Array(&ids), &contest.StartDate, &contest.Duration)
	if err == nil {
		contest.ProblemList = int64toInt(ids)
	}
	return contest, err
}