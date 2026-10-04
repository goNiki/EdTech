-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS headline VARCHAR(255);

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS headline;
