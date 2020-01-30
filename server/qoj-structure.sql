CREATE TABLE users (
	username CHARACTER(16),
	password CHARACTER(60),
	fullname CHARACTER(100),
	primary key (username)
);

CREATE TABLE problems (
	id	    SERIAL,
	code    CHAR(10),
	name    CHAR(100),
	tl      FLOAT,
	ml      INT,
	primary key (id)
);

CREATE TABLE submissions (
    id          SERIAL,
    username    CHARACTER(16) REFERENCES users(username),
    problem_id  INT REFERENCES problems(id),
    created_at  TIMESTAMP DEFAULT NOW(),
    status		TEXT,
    primary key (id)
);

CREATE TABLE tests (
	id			SERIAL,
	problem_id	INT REFERENCES problems(id),
	ord			INT,
	inp_preview	TEXT,
	out_preview	TEXT,
	primary key (id)
);

CREATE TABLE submissions (
    id          SERIAL,
    username    CHARACTER(16) REFERENCES users(username),
    problem_id  INT REFERENCES problems(id),
    created_at  TIMESTAMP DEFAULT NOW(),
    status		TEXT DEFAULT 'In queue...',
    primary key (id)
);