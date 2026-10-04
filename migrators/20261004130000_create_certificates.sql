-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS certificates (
    id BIGSERIAL PRIMARY KEY,
    certificate_code VARCHAR(64) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    student_name TEXT NOT NULL,
    course_title TEXT NOT NULL,
    final_score NUMERIC(5, 2) NOT NULL DEFAULT 0.00,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_user_course_certificate UNIQUE (user_id, course_id)
);

CREATE INDEX IF NOT EXISTS idx_certificates_code ON certificates(certificate_code);
CREATE INDEX IF NOT EXISTS idx_certificates_user ON certificates(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS certificates CASCADE;
-- +goose StatementEnd
