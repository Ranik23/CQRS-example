-- +goose Up
-- +goose StatementBegin

CREATE TYPE outbox_status AS ENUM ('not sent', 'sent', 'processing');

CREATE TABLE outbox (
    id SERIAL PRIMARY KEY,
    status outbox_status NOT NULL DEFAULT 'not sent',
    key     TEXT NOT NULL,
    message BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox;
DROP TYPE IF EXISTS outbox_status;
-- +goose StatementEnd
