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
select count(*) from users where username = ?