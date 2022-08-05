-- name: CreateTransfer :execresult
INSERT INTO transfers(
    from_account_id,
    to_account_id,
    amount
) values (
    ?,?,?
);

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = ?;

-- name: ListTransfer :many
SELECT * FROM transfers;

-- name: UpdateTransfer :execresult
UPDATE transfers
SET amount = ?
WHERE id = ?;

-- name: DeleteTransfer :execresult
DELETE FROM transfers
WHERE id = ?;