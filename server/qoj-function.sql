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