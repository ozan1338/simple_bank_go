-- name: CreateAccount :execresult
INSERT INTO accounts (
    owner,
    balance,
    currency
) VALUES (
    ?,?,?
);

-- name: GetLastInsertId :one
SELECT LAST_INSERT_ID();

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = ? LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts WHERE id = ? LIMIT 1 FOR UPDATE;

-- name: ListAccount :many
SELECT * FROM accounts
ORDER BY id
LIMIT ?
OFFSET ?;

-- name: UpdateAccount :execresult
UPDATE accounts
SET balance = ?
WHERE id = ?;

-- name: AddAccountBalance :execresult
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id);

-- name: DeleteAccount :execresult
DELETE FROM accounts
WHERE id = ?;