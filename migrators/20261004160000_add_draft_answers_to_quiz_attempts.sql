-- +goose Up
ALTER TABLE quiz_attempts ADD COLUMN IF NOT EXISTS draft_answers JSONB DEFAULT '{}'::jsonb;
ALTER TABLE quiz_attempts ADD COLUMN IF NOT EXISTS current_step INT DEFAULT 1;

-- +goose Down
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS current_step;
ALTER TABLE quiz_attempts DROP COLUMN IF EXISTS draft_answers;
