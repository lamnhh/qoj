CREATE OR REPLACE FUNCTION create_user(
	_username CHARACTER(16),
	_password CHARACTER(60),
	_fullname CHARACTER(100)
)
RETURNS SETOF users AS $$
BEGIN
	-- Verify unique username
	IF EXISTS (SELECT * FROM users WHERE username = _username) THEN
		RAISE unique_violation USING HINT = 'Username "' || _username || '" has been used';
	END IF;

	-- Insert new user row into table `users`
	RETURN QUERY
	INSERT INTO
		users
	VALUES
		(_username, _password, _fullname)
	RETURNING
		*;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_submission_result(sub_id INT)
RETURNS TABLE (
	inp_preview		TEXT,
	out_preview		TEXT,
	answer_preview	TEXT,
	score			FLOAT,
	verdict			TEXT
)
AS $$
BEGIN
	RETURN QUERY
	SELECT
		tests.inp_preview,
		tests.out_preview,
		submission_results.answer_preview,
		submission_results.score,
		submission_results.verdict
	FROM
		submission_results
		JOIN tests ON (submission_results.test_id = tests.id)
	ORDER BY
		tests.ord ASC;
END;
$$ LANGUAGE plpgsql;