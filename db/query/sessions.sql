-- name: CreateRefreshToken :execresult
INSERT INTO sessions (
    id,
    username,
    refresh_token,
    user_agent,
    client_ip,
    is_blocked,
    expired_at
) VALUES (
    ?,?,?,?,?,?,?
);

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = ? LIMIT 1;