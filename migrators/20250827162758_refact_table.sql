-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id SERIAL PRIMARY KEY, 
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL, 
    password_hash VARCHAR(255) NOT NULL, 
    role VARCHAR(50) NOT NULL DEFAULT 'user', 
    create_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, 
    update_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE users(
    id SERIAL PRIMARY KEY, 
    email VARCHAR(255) UNIQUE NOT NULL, 
    password_hash VARCHAR(255) NOT NULL, 
    role VARCHAR(50) NOT NULL DEFAULT 'user', 
    create_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, 
    update_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 
-- +goose StatementEnd
