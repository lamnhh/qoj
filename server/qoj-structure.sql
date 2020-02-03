CREATE TABLE users (
	username            CHARACTER(16),
	password            CHARACTER(60),
	fullname            CHARACTER(100),
	profile_picture     TEXT DEFAULT '/static/profile-picture-placeholder.png',
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
    username    CHARACTER(16),
    problem_id  INT,
    created_at  TIMESTAMP DEFAULT NOW(),
    status		TEXT DEFAULT 'In queue...',
    primary key (id)
);

CREATE TABLE tests (
	id			SERIAL,
	problem_id	INT,
	ord			INT,
	inp_preview	TEXT,
	out_preview	TEXT,
	primary key (id)
);

CREATE TABLE submission_results (
    submission_id   INT,
    test_id         INT,
    score           FLOAT,
    verdict         TEXT,
    answer_preview  TEXT,
    primary key (submission_id, test_id)
);

ALTER TABLE submissions
	ADD FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE,
	ADD FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE;

ALTER TABLE tests
	ADD FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE;

ALTER TABLE submission_results
	ADD FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE,
	ADD FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE;