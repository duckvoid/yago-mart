-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id BIGINT PRIMARY KEY,
    user_name TEXT REFERENCES users(name),
    status TEXT,
    accrual NUMERIC DEFAULT 0,
    created_date TIMESTAMP NOT NULL DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
