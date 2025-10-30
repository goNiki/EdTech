-- +goose Up
-- +goose StatementBegin
CREATE TYPE resource_type AS ENUM ('file', 'url');

CREATE TABLE resources (
    id BIGSERIAL PRIMARY KEY, 
    lesson_id BIGINT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE, 
    type resource_type NOT NULL, 
    path TEXT NOT NULL, 
    mime TEXT, 
    size BIGINT, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resources;
DROP TYPE IF EXISTS resource_type;
-- +goose StatementEnd
