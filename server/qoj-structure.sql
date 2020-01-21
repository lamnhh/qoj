CREATE TABLE users (
	username CHARACTER(16),
	password CHARACTER(60),
	fullname CHARACTER(100),
	primary key (username)
);

CREATE TABLE problems (
	id	SERIAL,
	code CHAR(10),
	name CHAR(10),
	primary key (id)
);