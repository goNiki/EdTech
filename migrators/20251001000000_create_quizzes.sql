-- +goose Up
CREATE TYPE quiz_type AS ENUM (
    'single_choice',      -- Один правильный ответ
    'multiple_choice',    -- Несколько правильных
    'true_false',         -- Истина/Ложь
    'fill_blank',         -- Заполнить пропуск
    'short_answer',       -- Короткий ответ
    'essay'               -- Эссе (открытый ответ)
);

CREATE TABLE quizzes (
    id BIGSERIAL PRIMARY KEY,
    
    lesson_id BIGINT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    
    title TEXT NOT NULL,
    description TEXT,
    
    -- Настройки
    type quiz_type NOT NULL,
    points INT DEFAULT 1,          -- Баллы за правильный ответ
    time_limit INT,                -- Лимит времени (секунды)
    max_attempts INT,              -- Макс. попыток (NULL = unlimited)
    passing_score INT DEFAULT 70,  -- Проходной балл (%)
    
    -- Сортировка
    position INT NOT NULL DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_quizzes_lesson ON quizzes(lesson_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_quizzes_course ON quizzes(course_id) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS quizzes CASCADE;
DROP TYPE IF EXISTS quiz_type;
