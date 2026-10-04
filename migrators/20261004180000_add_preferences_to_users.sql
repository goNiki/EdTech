-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN IF NOT EXISTS preferences JSONB DEFAULT '{"font_scale": "medium", "content_width": "wide", "line_height": "normal", "reading_theme": "system"}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS preferences;
-- +goose StatementEnd
