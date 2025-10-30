-- +goose Up
-- +goose StatementBegin
    CREATE TABLE lessons (
        id BIGSERIAL PRIMARY KEY, 
        course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE, 
        title TEXT NOT NULL, 
        description TEXT NOT NULL,
        cover_url TEXT,
        content TEXT, 
        position INT NOT NULL, 
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS lessons; 

-- +goose StatementEnd
