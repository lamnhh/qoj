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

CREATE OR REPLACE FUNCTION get_problem_list(_username CHARACTER(16))
RETURNS TABLE (
	id			INT,
	code		CHARACTER(10),
	name		CHARACTER(100),
	tl			FLOAT,
	ml			INT,
	max_score	FLOAT
)
AS $$
BEGIN
	RETURN QUERY
	SELECT
		s.id, s.code, s.name, s.tl, s.ml,
		COALESCE(MAX(s.score), 0) as max_score
	FROM (
		SELECT
			problems.id,
			problems.code,
			problems.name,
			problems.tl,
			problems.ml,
			submissions.id as sid,
			submissions.username,
			SUM(CASE
				WHEN submissions.username = _username THEN submission_results.score
				ELSE 0
			END) as score
		FROM
			problems
			LEFT JOIN submissions ON (problems.id = submissions.problem_id)
			LEFT JOIN submission_results ON (submissions.id = submission_results.submission_id)
		GROUP BY
			problems.id, problems.code, problems.name, problems.tl, problems.ml,submissions.id, submissions.username) s 
	GROUP BY
		s.id, s.code, s.name, s.tl, s.ml
	ORDER BY
		s.id ASC;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_problem_by_id(_problem_id INT, _username CHARACTER(16))
RETURNS TABLE (
	id			INT,
	code		CHARACTER(10),
	name		CHARACTER(100),
	tl			FLOAT,
	ml			INT,
	max_score	FLOAT
)
AS $$
BEGIN
	RETURN QUERY
	SELECT
		s.id, s.code, s.name, s.tl, s.ml,
		COALESCE(MAX(s.score), 0) as max_score
	FROM (
		SELECT
			problems.id,
			problems.code,
			problems.name,
			problems.tl,
			problems.ml,
			submissions.id as sid,
			submissions.username,
			SUM(CASE
				WHEN submissions.username = _username THEN submission_results.score
				ELSE 0
			END) as score
		FROM
			problems
			LEFT JOIN submissions ON (problems.id = submissions.problem_id)
			LEFT JOIN submission_results ON (submissions.id = submission_results.submission_id)
		WHERE
			problems.id = _problem_id
		GROUP BY
			problems.id, problems.code, problems.name, problems.tl, problems.ml,submissions.id, submissions.username) s 
	GROUP BY
		s.id, s.code, s.name, s.tl, s.ml
	ORDER BY
		s.id ASC;
END;
$$ LANGUAGE plpgsql;