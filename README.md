# Exinity Assessment

> docker compose up --build

# Usage

## Withdraw
> curl --location 'localhost/transactions' \
--header 'Content-Type: application/json' \
--data '{
"currency": "USD",
"amount": 5.00,
"gateway": "gatewayA",
"type": "deposit"
}'

## Deposit
> curl --location 'localhost/transactions' \
--header 'Content-Type: application/json' \
--data '{
"currency": "EUR",
"amount": 10.00,
"gateway": "gatewayA",
"type": "withdraw"
}'

## Callback
> curl --location 'localhost/callbacks' \
--header 'Content-Type: application/json' \
--data '{
"gateway": "gatewayA",
"transaction_id": "42411ba7-9127-4ba0-8b41-896d4fbae897",
"status": "success"
}'

## Tests
> go test ./...