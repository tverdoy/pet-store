DROP TABLE IF EXISTS pets_tags;
DROP TABLE IF EXISTS photos;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS pets;

DROP TABLE IF EXISTS orders;

DROP TABLE IF EXISTS auth;
DROP TABLE IF EXISTS users;


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
);

CREATE TYPE PetStatus AS ENUM ('available', 'pending','sold');

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE pets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    status  PetStatus,

    category_id INTEGER REFERENCES categories (id)
);

CREATE TABLE pets_tags (
    pet_id INTEGER REFERENCES pets (id),
    tag_id INTEGER REFERENCES tags (id)
);

CREATE TABLE photos (
    id SERIAL PRIMARY KEY,
    pet_id INTEGER REFERENCES pets (id)
);

CREATE TYPE OrderStatus AS ENUM ('placed', 'approved','delivered');

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    pet_id INTEGER REFERENCES pets (id),
    ship_date TIMESTAMP,
    status OrderStatus,
    complete BOOLEAN DEFAULT false
);