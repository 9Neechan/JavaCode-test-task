-- name: CreateWallet :one
INSERT INTO wallet (
  balance
) VALUES (
  $1
) RETURNING *;

-- name: UpdateWalletBalance :one
UPDATE wallet
SET balance = balance + sqlc.arg(amount)
WHERE wallet_id = sqlc.arg(wallet_id)
RETURNING *;

-- name: GetWallet :one
SELECT * FROM wallet
WHERE wallet_id = $1 LIMIT 1;