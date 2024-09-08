-- +goose Up
CREATE TABLE transactions
(
    id              UUID PRIMARY KEY,
    type            VARCHAR(10)              NOT NULL,
    amount          NUMERIC(12, 2)           NOT NULL,
    currency        VARCHAR(3)               NOT NULL,
    gateway         VARCHAR(50)              NOT NULL,
    status          VARCHAR(20)              NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    external_txn_id VARCHAR(100)
);

-- +goose Down
DROP TABLE transactions;
