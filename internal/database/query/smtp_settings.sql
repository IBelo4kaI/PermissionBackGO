-- name: GetSMTPSettings :one
-- Синглтон: ровно одна строка, id фиксирован на стороне Go (settings.SingletonID).
SELECT
	id,
	host,
	port,
	username,
	password_encrypted,
	from_address,
	from_name,
	use_tls,
	updated_at
FROM
	smtp_settings
WHERE
	id = ?;

-- name: UpsertSMTPSettings :exec
INSERT INTO
	smtp_settings (id, host, port, username, password_encrypted, from_address, from_name, use_tls)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY
UPDATE host =
VALUES
	(host),
	port =
VALUES
	(port),
	username =
VALUES
	(username),
	password_encrypted =
VALUES
	(password_encrypted),
	from_address =
VALUES
	(from_address),
	from_name =
VALUES
	(from_name),
	use_tls =
VALUES
	(use_tls);
