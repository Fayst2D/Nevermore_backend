-- +goose Up
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(50) NOT NULL UNIQUE,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       phone_number VARCHAR(255),
                       password VARCHAR(255) NOT NULL,
                       role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'moderator', 'user')),
                       photo VARCHAR(255),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- +goose DOWN
DROP TABLE users;
