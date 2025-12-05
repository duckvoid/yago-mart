-- +goose Up
-- +goose StatementBegin
CREATE TABLE balance (
     user_id   INTEGER PRIMARY KEY REFERENCES users(id),
     user_name TEXT REFERENCES users(name),
     current   NUMERIC NOT NULL DEFAULT 0,
     withdrawn NUMERIC NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS balance;
-- +goose StatementEnd
