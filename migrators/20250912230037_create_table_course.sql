-- +goose Up
CREATE TYPE course_status AS ENUM ('draft', 'published', 'archived');
CREATE TYPE course_visibility AS ENUM ('public', 'private');
CREATE TYPE course_difficulty AS ENUM ('beginner', 'intermediate', 'advanced');

CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    short_description TEXT,
    description TEXT,
    
    cover_url TEXT,
    intro_video_url TEXT,
    
    created_by BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    
    visibility course_visibility NOT NULL DEFAULT 'public',
    status course_status NOT NULL DEFAULT 'draft',
    difficulty course_difficulty,
    
    language VARCHAR(10) DEFAULT 'ru',
    estimated_duration INT,
    category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    
    total_lessons INT DEFAULT 0,
    total_sections INT DEFAULT 0,
    enrolled_count INT DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_courses_slug ON courses(slug) WHERE deleted_at IS NULL;
CREATE INDEX idx_courses_status_visibility ON courses(status, visibility) WHERE deleted_at IS NULL;
CREATE INDEX idx_courses_created_by ON courses(created_by) WHERE deleted_at IS NULL;
CREATE INDEX idx_courses_published_at ON courses(published_at DESC) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS courses CASCADE;
DROP TYPE IF EXISTS course_difficulty;
DROP TYPE IF EXISTS course_visibility;
DROP TYPE IF EXISTS course_status;
