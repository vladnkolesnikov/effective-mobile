-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_subscriptions
(
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(50)   NOT NULL,
    price           DECIMAL(5, 0) NOT NULL,
    start_date      DATE          NOT NULL,
    expiration_date DATE,
    user_id         UUID REFERENCES users (id) ON DELETE CASCADE,
    created_at      timestamptz      DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamptz      DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_exp_date CHECK (
        start_date < expiration_date
        )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_subscriptions;
-- +goose StatementEnd