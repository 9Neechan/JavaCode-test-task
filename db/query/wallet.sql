-- name: CreateWallet :one
INSERT INTO wallet (
  balance
) VALUES (
  $1
) RETURNING *;

-- name: UpdateWalletBalance :one
UPDATE wallet
SET balance = balance + sqlc.arg(amount)
WHERE wallet_uuid = sqlc.arg(wallet_uuid)
RETURNING *;

-- name: GetWallet :one
SELECT * FROM wallet
WHERE wallet_uuid = $1 LIMIT 1;