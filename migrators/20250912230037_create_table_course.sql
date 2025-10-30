-- +goose Up
-- +goose StatementBegin
CREATE TYPE course_status AS ENUM ('draft', 'published');
CREATE TYPE course_visibility AS ENUM('public', 'private')

CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY, 
    title TEXT NOT NULL, 
    slug TEXT NOT NULL UNIQUE, 
    description TEXT, 
    cover_url TEXT, 
    created_by BIGINT NOT NULL, 
    visibility course_visibility NOT NULL DEFAULT 'public',
    status course_status NOT NULL DEFAULT 'draft', 
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS courses;
DROP TYPE IF EXISTS course_visibility;
DROP TYPE IF EXISTS course_status;
-- +goose StatementEnd
