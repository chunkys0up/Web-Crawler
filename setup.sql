DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL
);

ALTER TABLE users ADD CONSTRAINT unique_username UNIQUE (username);

INSERT INTO users (username) VALUES ('andrew');

SELECT * FROM users;


