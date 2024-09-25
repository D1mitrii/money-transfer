-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
    uuid uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    account_name varchar(255) NOT NULL,
    balance bigint NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accounts;
-- +goose StatementEnd