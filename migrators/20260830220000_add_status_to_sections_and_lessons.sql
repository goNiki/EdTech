-- +goose Up
ALTER TABLE sections ADD COLUMN IF NOT EXISTS status VARCHAR(50) NOT NULL DEFAULT 'draft';
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS status VARCHAR(50) NOT NULL DEFAULT 'draft';

CREATE INDEX IF NOT EXISTS idx_sections_status ON sections(course_id, status) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_lessons_status ON lessons(course_id, status) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_lessons_status;
DROP INDEX IF EXISTS idx_sections_status;
ALTER TABLE lessons DROP COLUMN IF EXISTS status;
ALTER TABLE sections DROP COLUMN IF EXISTS status;
