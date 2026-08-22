-- +goose Up
CREATE TABLE quiz_questions (
    id BIGSERIAL PRIMARY KEY,
    
    quiz_id BIGINT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    
    question_text TEXT NOT NULL,
    explanation TEXT,  -- Объяснение правильного ответа
    
    position INT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_quiz_questions_quiz ON quiz_questions(quiz_id);

-- +goose Down
DROP TABLE IF EXISTS quiz_questions CASCADE;
