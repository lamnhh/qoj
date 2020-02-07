package contest

import (
	"database/sql"
	"github.com/lib/pq"
)

func int64toInt(a []int64) []int {
	b := make([]int, 0)
	for _, x := range a {
		b = append(b, int(x))
	}
	return b
}

func parseContestFromSqlRow(row *sql.Row) (Contest, error) {
	contest := Contest{}
	ids := make([]int64, 0)
	err := row.Scan(&contest.Id, &contest.Name, pq.Array(&ids), &contest.StartDate, &contest.Duration)
	if err == nil {
		contest.ProblemList = int64toInt(ids)
	}
	return contest, err
}

func parseContestFromSqlRows(rows *sql.Rows) (Contest, error) {
	contest := Contest{}
	ids := make([]int64, 0)
	err := rows.Scan(&contest.Id, &contest.Name, pq.Array(&ids), &contest.StartDate, &contest.Duration)
	if err == nil {
		contest.ProblemList = int64toInt(ids)
	}
	return contest, err
}