-- +goose Up
CREATE TYPE lesson_type AS ENUM ('lecture', 'exercise', 'quiz', 'assignment');

CREATE TABLE lessons (
    id BIGSERIAL PRIMARY KEY, 
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE, 
    section_id BIGINT REFERENCES sections(id) ON DELETE CASCADE,
    title TEXT NOT NULL, 
    description TEXT NOT NULL,
    cover_url TEXT,
    content TEXT, 
    type lesson_type NOT NULL DEFAULT 'lecture',
    position INT NOT NULL, 
    duration INT,
    is_free BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_lessons_course_position ON lessons(course_id, position) WHERE deleted_at IS NULL;
CREATE INDEX idx_lessons_type ON lessons(type) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS lessons CASCADE; 
DROP TYPE IF EXISTS lesson_type;
