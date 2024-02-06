DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS auth;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(255),
    password VARCHAR(255),
    user_status int
);

CREATE TABLE auth (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id)
)