-- name: GetSessionByTokenHash :one
-- token хэшируется SHA-256 на стороне Go перед вызовом запроса.
SELECT id, user_id, token_hash, expires_at
FROM sessions
WHERE token_hash = ?
  AND expires_at > UTC_TIMESTAMP();

-- name: CreateSession :exec
INSERT INTO sessions (id, user_id, token_hash, expires_at)
VALUES (?, ?, ?, ?);

-- name: DeleteSessionsByUserID :exec
-- Нужно вызывать перед удалением пользователя (FK sessions.user_id -> users.id).
DELETE FROM sessions
WHERE user_id = ?;

-- name: DeleteSessionByID :exec
DELETE FROM sessions
WHERE id = ?;
