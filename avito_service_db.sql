DROP TABLE relations;
DROP TABLE segments;
DROP TABLE users;


CREATE TABLE users (
	id INT UNIQUE NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE segments (
	id SERIAL PRIMARY KEY,
	name VARCHAR(20) NOT NULL
);

CREATE TABLE relations (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	segment_id INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (segment_id) REFERENCES segments(id)
);
