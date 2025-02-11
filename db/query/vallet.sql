-- name: CreateVallet :one
INSERT INTO vallet (
  balance
) VALUES (
  $1
) RETURNING *;

-- name: UpdateValletBalance :one
UPDATE vallet
SET balance = balance + sqlc.arg(amount)
WHERE vallet_id = sqlc.arg(vallet_id)
RETURNING *;

-- name: GetVallet :one
SELECT * FROM vallet
WHERE vallet_id = $1 LIMIT 1;