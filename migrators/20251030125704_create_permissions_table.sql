-- +goose Up
-- +goose StatementBegin

CREATE TYPE permission_action AS ENUM ('view', 'edit', 'delete', 'enroll', 'publish');
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    resource VARCHAR(50) NOT NULL,
    action permission_action NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(resource, action)
);

INSERT INTO permissions (resource, action, description) VALUES 
('course', 'view', 'view course'),
('course', 'edit', 'edit course'),
('course', 'delete', 'delete course'),
('course', 'enroll', 'enroll in course'),
('course', 'publish', 'publish course'),
('lesson', 'view', 'view lesson'),
('lesson', 'edit', 'edit lesson'),
('lesson', 'delete', 'delete lesson'),
('lesson', 'publish', 'publish lesson'),
('resource', 'view', 'view resource'),
('resource', 'edit', 'edit resource'),
('resource', 'delete', 'delete resource'),
('user', 'view', 'view user'),
('user', 'edit', 'edit user'),
('user', 'delete', 'delete user');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS permissions;
DROP TYPE IF EXISTS permission_action;
-- +goose StatementEnd
