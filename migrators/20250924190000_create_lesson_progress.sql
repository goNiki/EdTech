-- +goose Up
CREATE TYPE progress_status AS ENUM ('not_started', 'in_progress', 'completed');

CREATE TABLE lesson_progress (
    id BIGSERIAL PRIMARY KEY,
    
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id BIGINT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,  -- Денормализация
    
    status progress_status NOT NULL DEFAULT 'not_started',
    
    -- Для видео
    last_position INT DEFAULT 0,  -- Секунда, на которой остановился
    watch_time INT DEFAULT 0,     -- Общее время просмотра (секунды)
    
    -- Для упражнений/тестов
    score DECIMAL(5,2),           -- Балл (0-100)
    attempts INT DEFAULT 0,       -- Количество попыток
    
    -- Timestamps
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    last_accessed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    UNIQUE(user_id, lesson_id)
);

-- Индексы
CREATE INDEX idx_lesson_progress_user_course ON lesson_progress(user_id, course_id);
CREATE INDEX idx_lesson_progress_status ON lesson_progress(status);
CREATE INDEX idx_lesson_progress_completed ON lesson_progress(user_id, course_id, completed_at);

-- +goose Down
DROP TABLE IF EXISTS lesson_progress CASCADE;
DROP TYPE IF EXISTS progress_status;
