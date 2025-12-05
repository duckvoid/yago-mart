-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id BIGINT,
    user_id   INTEGER PRIMARY KEY REFERENCES users(id),
    user_name TEXT REFERENCES users(name),
    status TEXT,
    accrual INT DEFAULT 0,
    created_date TIMESTAMP NOT NULL DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
