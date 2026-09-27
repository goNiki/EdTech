-- +goose Up
ALTER TABLE quiz_attempt_answers ADD COLUMN IF NOT EXISTS points INT DEFAULT 0;
ALTER TABLE quiz_attempt_answers ADD COLUMN IF NOT EXISTS feedback TEXT;
ALTER TABLE quiz_attempt_answers ADD COLUMN IF NOT EXISTS attachment_url TEXT;
ALTER TABLE quiz_attempt_answers ADD COLUMN IF NOT EXISTS graded_by BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE quiz_attempt_answers ADD COLUMN IF NOT EXISTS graded_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_quiz_attempt_answers_grading ON quiz_attempt_answers(attempt_id) WHERE is_correct IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_quiz_attempt_answers_grading;
ALTER TABLE quiz_attempt_answers DROP COLUMN IF EXISTS graded_at;
ALTER TABLE quiz_attempt_answers DROP COLUMN IF EXISTS graded_by;
ALTER TABLE quiz_attempt_answers DROP COLUMN IF EXISTS attachment_url;
ALTER TABLE quiz_attempt_answers DROP COLUMN IF EXISTS feedback;
ALTER TABLE quiz_attempt_answers DROP COLUMN IF EXISTS points;
