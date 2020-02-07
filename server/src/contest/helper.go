package contest

import (
	"database/sql"
	"github.com/lib/pq"
	"qoj/server/src/problem"
)

func parseSingleContest(row *sql.Row, username string) (SingleContest, error) {
	contest := SingleContest{}
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

func parseMultipleContests(rows *sql.Rows) (MultipleContest, error) {
	contest := MultipleContest{}
	err := rows.Scan(
		&contest.Id,
		&contest.Name,
		&contest.StartDate,
		&contest.Duration,
		&contest.NumberOfParticipants,
		&contest.IsRegistered,
	)
	return contest, err
}