package contest

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"qoj/server/config"
	"time"
)

type Contest struct {
	Id          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	ProblemList []int     `json:"problemList" binding:"required"`
	StartDate   time.Time `json:"startDate" binding:"required"`
	Duration    int       `json:"duration" binding:"required"`
}

// createContest receives a `contest` which contains name, problemList, startDate and duration
// then creates an entry in the database for it and returns a new Contest containing the newly created ID
func createContest(contest Contest) (Contest, error) {
	cmd := `SELECT * FROM create_contest($1, $2, $3, $4)`
	err := config.DB.
		QueryRow(cmd, contest.Name, pq.Array(contest.ProblemList), contest.StartDate, contest.Duration).
		Scan(&contest.Id)
	return contest, err
}

func fetchAllContests() ([]Contest, error) {
	cmd := `
	SELECT 
		contests.id,
		RTRIM(contests.name),
		array_agg(problems.id) as problem_list,
		contests.start_date,
		contests.duration
	FROM
		contests
		JOIN problems ON contests.id = problems.contest_id
	GROUP BY
		contests.id,
		contests.name,
		contests.start_date,
		contests.duration`
	rows, err := config.DB.Query(cmd)
	if err != nil {
		return nil, err
	}

	contestList := make([]Contest, 0)
	for rows.Next() {
		contest, err := parseContestFromSqlRows(rows)
		if err == nil {
			contestList = append(contestList, contest)
		}
	}

	return contestList, nil
}

func fetchContestById(contestId int) (Contest, error) {
	cmd := `
	SELECT 
		contests.id,
		RTRIM(contests.name),
		array_agg(problems.id) as problem_list,
		contests.start_date,
		contests.duration
	FROM
		contests
		JOIN problems ON contests.id = problems.contest_id
	WHERE
		contests.id = $1
	GROUP BY
		contests.id,
		contests.name,
		contests.start_date,
		contests.duration`

	contest, err := parseContestFromSqlRow(config.DB.QueryRow(cmd, contestId))
	if err == sql.ErrNoRows {
		err = errors.New(fmt.Sprintf("Contest #%d does not exist", contestId))
	}
	return contest, err
}

func joinContest(contestId int, username string) error {
	_, err := config.DB.Exec("INSERT INTO contest_registrations VALUES ($1, $2)", contestId, username)
	return err
}

func fetchParticipantList(contestId int) ([]string, error) {
	cmd := "SELECT RTRIM(username) FROM contest_registrations WHERE contest_id = $1 ORDER BY username;"
	rows, err := config.DB.Query(cmd, contestId)
	if err != nil {
		return []string{}, err
	}

	participantList := make([]string, 0)
	for rows.Next() {
		participant := ""
		if err := rows.Scan(&participant); err == nil {
			participantList = append(participantList, participant)
		}
	}
	return participantList, nil
}