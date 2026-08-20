-- +goose Up

ALTER TABLE users
ADD CONSTRAINT users_email_key UNIQUE (email),
ADD COLUMN created_at timestamptz NOT NULL DEFAULT now();

-- +goose Down

ALTER TABLE users
DROP COLUMN created_at,
DROP CONSTRAINT users_email_key;
