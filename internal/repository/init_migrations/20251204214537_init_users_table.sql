-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
     id     SERIAL PRIMARY KEY,
     name  TEXT NOT NULL UNIQUE,
     password TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
