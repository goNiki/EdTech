-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_role AS ENUM ('student', 'teacher', 'creator');

CREATE TABLE users_courses (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE, 
    role user_role NOT NULL DEFAULT 'student',
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, course_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_courses;
DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd
