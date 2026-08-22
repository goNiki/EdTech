-- +goose Up
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon_url TEXT,
    
    -- Иерархия (для подкатегорий)
    parent_id BIGINT REFERENCES categories(id) ON DELETE CASCADE,
    
    -- Сортировка
    position INT NOT NULL DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_slug ON categories(slug);

-- Примеры данных
INSERT INTO categories (name, slug, description, position) VALUES
('Программирование', 'programming', 'Курсы по программированию', 1),
('Веб-разработка', 'web-development', 'Frontend и Backend разработка', 2),
('Data Science', 'data-science', 'Наука о данных и ML', 3),
('Дизайн', 'design', 'UI/UX и графический дизайн', 4),
('Бизнес', 'business', 'Бизнес и менеджмент', 5);

-- +goose Down
DROP TABLE IF EXISTS categories CASCADE;
