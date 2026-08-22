-- +goose Up
CREATE TABLE quiz_answers (
    id BIGSERIAL PRIMARY KEY,
    
    question_id BIGINT NOT NULL REFERENCES quiz_questions(id) ON DELETE CASCADE,
    
    answer_text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    
    position INT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_quiz_answers_question ON quiz_answers(question_id);

-- +goose Down
DROP TABLE IF EXISTS quiz_answers CASCADE;
