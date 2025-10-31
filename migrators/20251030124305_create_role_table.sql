-- +goose Up
-- +goose StatementBegin

CREATE TYPE role_type AS ENUM ('student', 'teacher', 'admin');

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type role_type NOT NULL DEFAULT 'student',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO roles (name, type) VALUES 
('student', 'student'),
('teacher', 'teacher'),
('admin', 'admin');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
DROP TYPE IF EXISTS role_type;
-- +goose StatementEnd
