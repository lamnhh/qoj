package contest

import (
	"database/sql"
)

func parseContestFromRow(row *sql.Row) (Contest, error) {
	contest := Contest{}
	err := row.Scan(
		&contest.Id,
		&contest.Name,
		&contest.StartDate,
		&contest.Duration,
		&contest.NumberOfParticipants,
		&contest.IsRegistered,
	)
	return contest, err
}

func parseContestFromRows(rows *sql.Rows) (Contest, error) {
	contest := Contest{}
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