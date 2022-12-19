DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id varchar(32) PRIMARY KEY,
    email varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    create_at TIMESTAMP NOT NULL DEFAULT now()
);

DROP TABLE IF EXISTS posts;

CREATE TABLE posts (
    id varchar(32) PRIMARY KEY,
    content varchar(255) NOT NULL,
    create_at TIMESTAMP NOT NULL DEFAULT now(),
    user_id varchar(32) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
)
