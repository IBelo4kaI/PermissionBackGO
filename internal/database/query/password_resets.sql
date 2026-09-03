-- name: CreatePasswordReset :exec
INSERT INTO
	password_resets (id, user_id, token_hash, expires_at)
VALUES
	(?, ?, ?, ?);

-- name: GetPasswordResetByTokenHash :one
-- token хэшируется SHA-256 на стороне Go перед вызовом запроса (как sessions).
SELECT
	id,
	user_id,
	token_hash,
	used,
	expires_at,
	created_at
FROM
	password_resets
WHERE
	token_hash = ?
	AND used = 0
	AND expires_at > UTC_TIMESTAMP();

-- name: MarkPasswordResetUsed :exec
UPDATE password_resets
SET
	used = 1
WHERE
	id = ?;

-- name: DeletePasswordResetsByUserID :exec
-- Нужно вызывать перед удалением пользователя (FK password_resets.user_id -> users.id), как DeleteSessionsByUserID.
DELETE FROM password_resets
WHERE
	user_id = ?;
