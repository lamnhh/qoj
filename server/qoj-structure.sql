CREATE TABLE users (
	username            CHARACTER(16),
	password            CHARACTER(60),
	fullname            CHARACTER(100),
	profile_picture     TEXT DEFAULT '/static/images/profile-picture-placeholder.png',
	primary key (username)
);

CREATE TABLE languages (
	id		SERIAL,
	name	CHARACTER(30),
	ext		CHARACTER(10),
	command	TEXT,
	primary key (id)
);

CREATE TABLE problems (
	id	        SERIAL,
	code        CHAR(10),
	name        CHAR(100),
	tl          FLOAT,
	ml          INT,
	contest_id  INT,
	original_id INT,
	primary key (id)
);

CREATE TABLE submissions (
    id          SERIAL,
    username    CHARACTER(16),
    problem_id  INT,
    created_at  TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'utc'),
    status		TEXT DEFAULT 'In queue...',
    code        TEXT,
    compile_msg TEXT,
    language_id INT,
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
    execution_time  FLOAT,
    memory_used     INT,
    primary key (submission_id, test_id)
);

CREATE TABLE contests (
	id			SERIAL,
	name		CHARACTER(100),
	start_date	TIMESTAMP,
	duration	INT, -- duration in minutes
	primary key (id)
);

CREATE TABLE contest_registrations (
	contest_id 	INT,
	username	CHARACTER(16),
	primary key (contest_id, username)
);

ALTER TABLE submissions
	ADD FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE,
	ADD FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE,
	ADD FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE CASCADE;

ALTER TABLE tests
	ADD FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE;

ALTER TABLE submission_results
	ADD FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE,
	ADD FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE;

ALTER TABLE problems
	ADD FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE;

ALTER TABLE contest_registrations
	ADD FOREIGN KEY (contest_id) REFERENCES contests(id) ON DELETE CASCADE,
	ADD FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE;

CREATE OR REPLACE FUNCTION insert_problem_original_id() RETURNS TRIGGER AS $$
BEGIN
	IF (NEW.original_id IS NULL) THEN
		NEW.original_id = NEW.id;
	END IF;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_problem_insert
BEFORE INSERT ON problems
FOR EACH ROW EXECUTE PROCEDURE insert_problem_original_id();

INSERT INTO languages(name, ext, command) VALUES
	('C', '.c', 'gcc -Wall -lm -static -DEVAL -o %s -O2 %s.c'),
	('C++', '.cpp', 'g++ -Wall -lm -static -DEVAL -o %s -O2 %s.cpp'),
	('C++11', '.cpp', 'g++-7 -Wall -lm -static -DEVAL -o %s -O2 %s.cpp -std=c++11'),
	('C++14', '.cpp', 'g++-7 -Wall -lm -static -DEVAL -o %s -O2 %s.cpp -std=c++14'),
	('C++17', '.cpp', 'g++-7 -Wall -lm -static -DEVAL -o %s -O2 %s.cpp -std=c++17');