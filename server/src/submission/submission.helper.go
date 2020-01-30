package submission

import "qoj/server/config"

func updateScore(submissionId int, testId int, score float32, verdict string, preview string) error {
	_, err := config.DB.Exec("INSERT INTO submission_results VALUES ($1, $2, $3, $4, $5)",
		submissionId,
		testId,
		score,
		verdict,
		preview,
	)
	return err
}
