-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2,$3
)RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts 
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts 
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;

-- name: SelectForUpdate :one
BEGIN;
SELECT * FROM accounts 
WHERE id = $1 LIMIT 1
FOR UPDATE;

-- name: SelectAccWithEntry :many
SELECT
  entries.amount AS amount,
  accounts.id AS account_id
FROM accounts
JOIN entries ON accounts.id=entries.account_id;


-- name: ListAccountsFilter :many
SELECT *
FROM accounts
WHERE (@currency::text = '' OR currency = @currency)
    AND (@created_at::date = '0001-01-01' OR created_at > @created_at )
    AND (@owner::text = '' OR owner ILIKE '%' || @owner || '%');