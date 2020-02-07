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
	verdict			TEXT,
	execution_time	FLOAT,
	memory_used		INT
)
AS $$
BEGIN
	RETURN QUERY
	SELECT
		tests.inp_preview,
		tests.out_preview,
		submission_results.answer_preview,
		submission_results.score,
		submission_results.verdict,
		submission_results.execution_time,
		submission_results.memory_used
	FROM
		submission_results
		JOIN tests ON (submission_results.test_id = tests.id)
	WHERE
	    submission_id = sub_id
	ORDER BY
		tests.ord ASC;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_problem_list(_username CHARACTER(16), _page INT, _size INT)
RETURNS TABLE (
	id			INT,
	code		CHARACTER(10),
	name		CHARACTER(100),
	tl			FLOAT,
	ml			INT,
	max_score	FLOAT,
	test_count	INT
)
AS $$
BEGIN
	RETURN QUERY
	SELECT
        s.id, s.code, s.name, s.tl, s.ml, s.max_score,
        COUNT(*)::int
    FROM (
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
                problems.id, problems.code, problems.name, problems.tl, problems.ml,submissions.id, submissions.username
			ORDER BY
				problems.id DESC) s
        GROUP BY
            s.id, s.code, s.name, s.tl, s.ml) s
        LEFT JOIN tests ON (s.id = tests.problem_id)
    GROUP BY
        s.id, s.code, s.name, s.tl, s.ml, s.max_score
	ORDER BY
		s.id DESC
	OFFSET _page * _size LIMIT _size;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_problem_by_id(_problem_id INT, _username CHARACTER(16))
RETURNS TABLE (
	id			INT,
	code		CHARACTER(10),
	name		CHARACTER(100),
	tl			FLOAT,
	ml			INT,
	max_score	FLOAT,
	test_count	INT
)
AS $$
BEGIN
	RETURN QUERY
	SELECT
        s.id, s.code, s.name, s.tl, s.ml, s.max_score,
        COUNT(*)::int
    FROM (
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
            s.id ASC
    ) s
        LEFT JOIN tests ON (s.id = tests.problem_id)
    GROUP BY
        s.id, s.code, s.name, s.tl, s.ml, s.max_score;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION create_contest(
	_name 			CHARACTER(100),
	_problem_ids 	INT[],
	_start_date 	TIMESTAMP,
	_duration		INT
)
RETURNS INT
AS $$
DECLARE
	_contest_id		INT;
BEGIN
	INSERT INTO
		contests(name, start_date, duration)
	VALUES
		(_name, _start_date, _duration)
	RETURNING
		id INTO _contest_id;

	INSERT INTO
		problems(code, name, tl, ml, original_id, contest_id)
	SELECT
		code,
		name,
		tl,
		ml,
		problems.id as original_id,
		_contest_id as contest_id
	FROM
		problems
	WHERE
		problems.id = ANY(_problem_ids);

	RETURN _contest_id;
END;
$$LANGUAGE plpgsql;