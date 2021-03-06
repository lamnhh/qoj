package contest

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"qoj/server/config"
	"qoj/server/src/problem"
	"strings"
	"time"
)

type Contest struct {
	Id                   int       `json:"id"`
	Name                 string    `json:"name" binding:"required"`
	StartDate            time.Time `json:"startDate" binding:"required"`
	ProblemList          []int     `json:"problemList,omitempty" binding:"required"`
	Duration             int       `json:"duration" binding:"required"`
	NumberOfParticipants int       `json:"numberOfParticipants"`
	IsRegistered         bool      `json:"isRegistered"`
}

type ContestScore struct {
	ProblemId int     `json:"problemId"`
	Score     float32 `json:"score"`
}

type ContestScoreList struct {
	Username  string         `json:"username"`
	ScoreSum  float32        `json:"scoreSum"`
	ScoreList []ContestScore `json:"scoreList"`
}

func (s *ContestScore) Scan(src interface{}) error {
	bs, ok := src.([]byte)
	if ok == false {
		return errors.New("Not a []byte")
	}
	return json.Unmarshal(bs, s)
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

func fetchAllContests(username string) ([]Contest, error) {
	cmd := `
	SELECT 
		contests.id,
		RTRIM(contests.name),
		contests.start_date,
		contests.duration,
		COUNT(username) as participants,
		MAX(CASE
			WHEN username = $1 THEN 1 ELSE 0
		END) as is_registered
	FROM
		contests
		LEFT JOIN contest_registrations ON (contests.id = contest_registrations.contest_id)
	GROUP BY
		contests.id,
		contests.name,
		contests.start_date,
		contests.duration
	ORDER BY
		contests.start_date DESC;`
	rows, err := config.DB.Query(cmd, username)
	if err != nil {
		return nil, err
	}

	contestList := make([]Contest, 0)
	for rows.Next() {
		contest, err := parseContestFromRows(rows)
		if err == nil {
			contestList = append(contestList, contest)
		}
	}

	return contestList, nil
}

func fetchContestById(contestId int, username string) (Contest, error) {
	cmd := `
	SELECT 
		contests.id,
		RTRIM(contests.name),
		contests.start_date,
		contests.duration,
		COUNT(username) as participants,
		MAX(CASE
			WHEN username = $2 THEN 1 ELSE 0
		END) as is_registered
	FROM
		contests
		LEFT JOIN contest_registrations ON (contests.id = contest_registrations.contest_id)
	WHERE
		contests.id = $1
	GROUP BY
		contests.id,
		contests.name,
		contests.start_date,
		contests.duration`

	contest, err := parseContestFromRow(config.DB.QueryRow(cmd, contestId, username))
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

func fetchProblemList(contestId int, username string) ([]problem.Problem, error) {
	rows, err := config.DB.Query("SELECT id FROM problems WHERE contest_id = $1 ORDER BY id ASC", contestId)
	if err != nil {
		return nil, err
	}

	ids := make([]int, 0)
	for rows.Next() {
		id := 0
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}

	return problem.FetchProblemByIds(ids, username)
}

func fetchContestScore(contestId int) ([]ContestScoreList, error) {
	rows, err := config.DB.Query("SELECT * FROM get_contest_scores($1)", contestId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("Contest with ID `%d` does not exist", contestId))
		}
		return nil, err
	}

	scoreList := make([]ContestScoreList, 0)
	for rows.Next() {
		score := ContestScoreList{}
		if err := rows.Scan(&score.Username, &score.ScoreSum, pq.Array(&score.ScoreList)); err == nil {
			score.Username = strings.TrimSpace(score.Username)
			scoreList = append(scoreList, score)
		}
	}

	return scoreList, nil
}

func fetchContestByIdAdmin(contestId int) (Contest, error) {
	cmd := `SELECT
		contests.id,
		RTRIM(contests.name),
		contests.start_date,
		contests.duration,
		ARRAY_AGG(problems.original_id) as problem_list
	FROM
		contests
		LEFT JOIN problems ON (contests.id = problems.contest_id)
	WHERE
		contests.id = $1
	GROUP BY
		contests.id,
		contests.name,
		contests.start_date,
		contests.duration`

	var pid []int64
	contest := Contest{}
	err := config.DB.QueryRow(cmd, contestId).
		Scan(&contest.Id, &contest.Name, &contest.StartDate, &contest.Duration, pq.Array(&pid))
	if err != nil {
		return contest, err
	}

	contest.ProblemList = make([]int, 0)
	for _, x := range pid {
		contest.ProblemList = append(contest.ProblemList, int(x))
	}

	return contest, nil
}

func updateContest(contestId int, patch Contest) error {
	_, err := config.DB.Exec("SELECT update_contest($1, $2, $3, $4, $5)",
		contestId,
		patch.Name,
		patch.StartDate,
		patch.Duration,
		pq.Array(patch.ProblemList),
	)
	return err
}

func deleteContest(contestId int) error {
	_, err := config.DB.Exec("DELETE FROM contests WHERE id = $1", contestId)
	return err
}