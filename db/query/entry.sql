-- name: CreateEntries :execresult
INSERT INTO entries(
    account_id,
    amount
) values (
    ?,?
);

-- name: GetEntries :one
SELECT * FROM entries
WHERE id = ?;

-- name: ListEntries :many
SELECT * FROM entries;

-- name: UpdateEntries :execresult
UPDATE entries
SET amount = ?
WHERE id = ?;

-- name: GetIdEntries :one
SELECT * FROM entries
LIMIT 1;

-- name: DeleteEntries :execresult
DELETE FROM entries
WHERE id = ?;