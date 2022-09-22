-- name: CreateUser :execresult
INSERT INTO users (
    username,
    password,
    full_name,
    email
) VALUES (
    ?,?,?,?
);

-- name: GetUser :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: UserExist :one
select exists(select * from users where username = ?) as isExist;

-- name: UserMoreThanOne :one
select count(*) from users where username = ?;

-- name: UpdateUser :execresult
UPDATE users
SET
    password = COALESCE(sqlc.narg(password), password),
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email)
WHERE
    username = sqlc.arg(username);