-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id SERIAL PRIMARY KEY, 
    email VARCHAR(255) UNIQUE NOT NULL, 
    password_hash VARCHAR(255) NOT NULL, 
    role VARCHAR(50) NOT NULL DEFAULT 'user', 
    create_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, 
    update_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd