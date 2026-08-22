-- +goose Up
CREATE TABLE quiz_attempt_answers (
    id BIGSERIAL PRIMARY KEY,
    
    attempt_id BIGINT NOT NULL REFERENCES quiz_attempts(id) ON DELETE CASCADE,
    question_id BIGINT NOT NULL REFERENCES quiz_questions(id) ON DELETE CASCADE,
    answer_id BIGINT REFERENCES quiz_answers(id) ON DELETE SET NULL,  -- Для choice вопросов
    
    user_answer TEXT,  -- Для fill_blank, short_answer, essay
    is_correct BOOLEAN,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_quiz_attempt_answers_attempt ON quiz_attempt_answers(attempt_id);

-- +goose Down
DROP TABLE IF EXISTS quiz_attempt_answers CASCADE;
