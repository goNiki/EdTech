-- +goose Up
CREATE TABLE sections (
    id BIGSERIAL PRIMARY KEY,
    
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    
    title TEXT NOT NULL,
    description TEXT,
    
    -- Сортировка
    position INT NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_sections_course_id ON sections(course_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_sections_position ON sections(course_id, position) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS sections CASCADE;
