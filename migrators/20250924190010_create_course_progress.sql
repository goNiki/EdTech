-- +goose Up
CREATE TABLE course_progress (
    id BIGSERIAL PRIMARY KEY,
    
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    
    -- Прогресс
    completed_lessons INT DEFAULT 0,
    total_lessons INT NOT NULL,
    progress_percentage INT DEFAULT 0 CHECK (progress_percentage >= 0 AND progress_percentage <= 100),
    
    -- Времязатраты
    total_watch_time INT DEFAULT 0,  -- Секунды
    
    -- Средний балл по тестам
    average_score DECIMAL(5,2),
    
    -- Timestamps
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_accessed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    
    UNIQUE(user_id, course_id)
);

-- Индексы
CREATE INDEX idx_course_progress_user ON course_progress(user_id);
CREATE INDEX idx_course_progress_course ON course_progress(course_id);
CREATE INDEX idx_course_progress_percentage ON course_progress(progress_percentage);

-- Триггер для автоматического расчета прогресса
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_course_progress()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE course_progress cp
    SET 
        completed_lessons = (
            SELECT COUNT(*) 
            FROM lesson_progress lp 
            WHERE lp.user_id = NEW.user_id 
              AND lp.course_id = NEW.course_id 
              AND lp.status = 'completed'
        ),
        progress_percentage = (
            SELECT ROUND((COUNT(*) * 100.0 / cp.total_lessons))
            FROM lesson_progress lp
            WHERE lp.user_id = NEW.user_id
              AND lp.course_id = NEW.course_id
              AND lp.status = 'completed'
        ),
        last_accessed_at = NOW()
    WHERE cp.user_id = NEW.user_id AND cp.course_id = NEW.course_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trigger_update_course_progress
    AFTER INSERT OR UPDATE ON lesson_progress
    FOR EACH ROW
    EXECUTE FUNCTION update_course_progress();

-- +goose Down
DROP TRIGGER IF EXISTS trigger_update_course_progress ON lesson_progress;
DROP FUNCTION IF EXISTS update_course_progress();
DROP TABLE IF EXISTS course_progress CASCADE;
