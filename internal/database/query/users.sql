-- name: GetUserByID :one
SELECT
	id,
	name,
	surname,
	patronymic,
	username,
	gender_id,
	birthday,
	password,
	status,
	created_at
FROM
	users
WHERE
	id = ?;

-- name: GetUserByUsername :one
SELECT
	id,
	name,
	surname,
	patronymic,
	username,
	gender_id,
	birthday,
	password,
	status,
	created_at
FROM
	users
WHERE
	username = ?;

-- name: ListUsers :many
SELECT
	id,
	name,
	surname,
	patronymic,
	username,
	gender_id,
	birthday,
	password,
	status,
	created_at
FROM
	users
ORDER BY
	created_at DESC
LIMIT
	?
OFFSET
	?;

-- name: CountUsers :one
SELECT
	COUNT(*) AS total
FROM
	users;

-- name: ListAllUsers :many
SELECT
	id,
	name,
	surname,
	patronymic,
	username,
	gender_id,
	birthday,
	password,
	status,
	created_at
FROM
	users
ORDER BY
	created_at DESC;

-- name: ListUsersByServiceID :many
SELECT DISTINCT
	u.id,
	u.name,
	u.surname,
	u.patronymic,
	u.username,
	u.gender_id,
	u.birthday,
	u.password,
	u.status,
	u.created_at
FROM
	users u
	JOIN user_roles ur ON ur.user_id = u.id
	JOIN roles r ON r.id = ur.role_id
WHERE
	r.service_id = ?
ORDER BY
	u.created_at DESC
LIMIT
	?
OFFSET
	?;

-- name: CountUsersByServiceID :one
SELECT
	COUNT(DISTINCT u.id) AS total
FROM
	users u
	JOIN user_roles ur ON ur.user_id = u.id
	JOIN roles r ON r.id = ur.role_id
WHERE
	r.service_id = ?;

-- name: CreateUser :exec
INSERT INTO
	users (id, name, surname, patronymic, username, gender_id, birthday, password, status, created_at)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?, NOW());

-- name: UpdateUser :exec
-- Частичное обновление: sqlc.narg-поля, для которых передан NULL, остаются без изменений.
UPDATE users
SET
	name = COALESCE(sqlc.narg ('name'), name),
	surname = COALESCE(sqlc.narg ('surname'), surname),
	patronymic = COALESCE(sqlc.narg ('patronymic'), patronymic),
	username = COALESCE(sqlc.narg ('username'), username),
	birthday = COALESCE(sqlc.narg ('birthday'), birthday),
	status = COALESCE(sqlc.narg ('status'), status),
	gender_id = COALESCE(sqlc.narg ('gender_id'), gender_id),
	password = COALESCE(sqlc.narg ('password'), password)
WHERE
	id = sqlc.arg ('id');

-- name: DeleteUser :exec
DELETE FROM users
WHERE
	id = ?;

-- name: AddRoleToUser :exec
INSERT INTO
	user_roles (user_id, role_id)
VALUES
	(?, ?);

-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE
	user_id = ?
	AND role_id = ?;

-- name: ListRolesForUser :many
SELECT
	r.id,
	r.service_id,
	r.name,
	r.description,
	r.is_global,
	r.created_at
FROM
	roles r
	JOIN user_roles ur ON ur.role_id = r.id
WHERE
	ur.user_id = ?;
