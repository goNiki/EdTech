-- +goose Up
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS quiz_settings JSONB DEFAULT '{
  "time_limit_minutes": 0,
  "question_time_limit_seconds": 0,
  "max_attempts": 0,
  "passing_score_percent": 70,
  "feedback_mode": "immediate",
  "shuffle_questions": false
}'::jsonb;

-- +goose Down
ALTER TABLE lessons DROP COLUMN IF EXISTS quiz_settings;
