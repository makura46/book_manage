CREATE TABLE users (
	name VARCHAR(64) NOT NULL,
	password CHAR(64) NOT NULL,
	session CHAR(64),
	PRIMARY KEY(name)
);
