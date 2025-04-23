-- +goose Up:
ALTER TABLE users
ADD hashed_password TEXT NOT NULL
ADD CONSTRAINT df_hashed_password
DEFAULT "unset" FOR hashed_passowrd;

-- +goose Down:
ALTER TABLE users
DROP COLUMN hashed_password;