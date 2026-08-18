-- name: CreateInvite :exec
INSERT INTO
	invites (id, code, email, person_id, company_id, department_id, position_id, created_by, expires_at)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetInviteByID :one
SELECT
	id,
	code,
	email,
	person_id,
	company_id,
	department_id,
	position_id,
	created_by,
	user_id,
	used,
	revoked,
	expires_at,
	created_at
FROM
	invites
WHERE
	id = ?;

-- name: GetInviteByCode :one
SELECT
	id,
	code,
	email,
	person_id,
	company_id,
	department_id,
	position_id,
	created_by,
	user_id,
	used,
	revoked,
	expires_at,
	created_at
FROM
	invites
WHERE
	code = ?;

-- name: GetActiveInviteByEmail :one
-- Используется при создании инвайта: ищет ещё не использованный, не
-- отозванный и не истёкший инвайт на этот email, чтобы отозвать его перед
-- выдачей нового (см. invite.Service.Create).
SELECT
	id,
	code,
	email,
	person_id,
	company_id,
	department_id,
	position_id,
	created_by,
	user_id,
	used,
	revoked,
	expires_at,
	created_at
FROM
	invites
WHERE
	email = ?
	AND used = 0
	AND revoked = 0
	AND expires_at > UTC_TIMESTAMP();

-- name: ListInvites :many
-- Без пагинации (см. ListAllUsers): фильтрация по search/company_id/
-- department_id/position_id и сортировка выполняются в Go, см.
-- invite.Service.List.
SELECT
	id,
	code,
	email,
	person_id,
	company_id,
	department_id,
	position_id,
	created_by,
	user_id,
	used,
	revoked,
	expires_at,
	created_at
FROM
	invites
ORDER BY
	created_at DESC;

-- name: RevokeInvite :exec
UPDATE invites
SET
	revoked = 1
WHERE
	id = ?;

-- name: MarkInviteUsed :exec
UPDATE invites
SET
	used = 1,
	user_id = ?
WHERE
	id = ?;
