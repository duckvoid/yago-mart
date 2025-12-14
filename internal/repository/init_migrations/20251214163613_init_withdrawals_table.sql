-- +goose Up
-- +goose StatementBegin
CREATE TABLE withdrawals (
        id INT PRIMARY KEY,
        user_name TEXT REFERENCES users(name),
        order_id BIGINT REFERENCES orders(id),
        sum NUMERIC NOT NULL DEFAULT 0,
        processed_at TIMESTAMP NOT NULL DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals;
-- +goose StatementEnd
