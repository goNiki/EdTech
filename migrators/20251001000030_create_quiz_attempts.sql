-- +goose Up
CREATE TABLE quiz_attempts (
    id BIGSERIAL PRIMARY KEY,
    
    quiz_id BIGINT NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Результаты
    score DECIMAL(5,2) NOT NULL,  -- Балл (0-100)
    passed BOOLEAN NOT NULL,       -- Прошел или нет
    
    -- Время
    time_spent INT,  -- Секунды
    
    -- Timestamps
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_quiz_attempts_user ON quiz_attempts(user_id);
CREATE INDEX idx_quiz_attempts_quiz_user ON quiz_attempts(quiz_id, user_id);
CREATE INDEX idx_quiz_attempts_score ON quiz_attempts(score DESC);

-- +goose Down
DROP TABLE IF EXISTS quiz_attempts CASCADE;
