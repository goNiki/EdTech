-- +goose Up
CREATE TYPE resource_type AS ENUM (
    'video', 'document', 'link', 'image', 'audio', 'code', 'archive'
);

CREATE TABLE resources (
    id BIGSERIAL PRIMARY KEY,
    lesson_id BIGINT REFERENCES lessons(id) ON DELETE CASCADE,
    course_id BIGINT REFERENCES courses(id) ON DELETE CASCADE,
    title TEXT,
    description TEXT,
    type resource_type NOT NULL,
    path TEXT NOT NULL,
    mime TEXT,
    size BIGINT,
    external_url TEXT,
    duration INT,
    order_position INT DEFAULT 0,
    is_required BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    CONSTRAINT resources_parent_check CHECK (
        (lesson_id IS NOT NULL AND course_id IS NULL) OR 
        (lesson_id IS NULL AND course_id IS NOT NULL)
    )
);

CREATE INDEX idx_resources_lesson_id ON resources(lesson_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_resources_course_id ON resources(course_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_resources_type ON resources(type) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS resources CASCADE;
DROP TYPE IF EXISTS resource_type;
