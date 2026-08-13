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
-- search: поиск по name/surname/patronymic/username. sort_by: name/surname/
-- patronymic/username/birthday/status/created_at (иначе created_at),
-- sort_dir: asc (иначе desc) — см. query.QuerySort.
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
	users u
WHERE
	CAST(sqlc.narg (search) AS char) IS NULL
	OR u.name LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	OR u.surname LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	OR u.patronymic LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	OR u.username LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
ORDER BY
	CASE
		WHEN CAST(sqlc.arg (sort_dir) AS char) = 'asc' THEN CASE CAST(sqlc.arg (sort_by) AS char)
			WHEN 'name' THEN CAST(u.name AS char)
			WHEN 'surname' THEN CAST(u.surname AS char)
			WHEN 'patronymic' THEN CAST(u.patronymic AS char)
			WHEN 'username' THEN CAST(u.username AS char)
			WHEN 'birthday' THEN CAST(u.birthday AS char)
			WHEN 'status' THEN CAST(u.status AS char)
			ELSE CAST(u.created_at AS char)
		END
	END ASC,
	CASE
		WHEN CAST(sqlc.arg (sort_dir) AS char) = 'desc' THEN CASE CAST(sqlc.arg (sort_by) AS char)
			WHEN 'name' THEN CAST(u.name AS char)
			WHEN 'surname' THEN CAST(u.surname AS char)
			WHEN 'patronymic' THEN CAST(u.patronymic AS char)
			WHEN 'username' THEN CAST(u.username AS char)
			WHEN 'birthday' THEN CAST(u.birthday AS char)
			WHEN 'status' THEN CAST(u.status AS char)
			ELSE CAST(u.created_at AS char)
		END
	END DESC
LIMIT
	?
OFFSET
	?;

-- name: CountUsers :one
SELECT
	COUNT(*) AS total
FROM
	users u
WHERE
	CAST(sqlc.narg (search) AS char) IS NULL
	OR u.name LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	OR u.surname LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	OR u.patronymic LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	OR u.username LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%');

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
	r.service_id = sqlc.arg (service_id)
	AND (
		CAST(sqlc.narg (search) AS char) IS NULL
		OR u.name LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
		OR u.surname LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
		OR u.patronymic LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
		OR u.username LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	)
ORDER BY
	CASE
		WHEN CAST(sqlc.arg (sort_dir) AS char) = 'asc' THEN CASE CAST(sqlc.arg (sort_by) AS char)
			WHEN 'name' THEN CAST(u.name AS char)
			WHEN 'surname' THEN CAST(u.surname AS char)
			WHEN 'patronymic' THEN CAST(u.patronymic AS char)
			WHEN 'username' THEN CAST(u.username AS char)
			WHEN 'birthday' THEN CAST(u.birthday AS char)
			WHEN 'status' THEN CAST(u.status AS char)
			ELSE CAST(u.created_at AS char)
		END
	END ASC,
	CASE
		WHEN CAST(sqlc.arg (sort_dir) AS char) = 'desc' THEN CASE CAST(sqlc.arg (sort_by) AS char)
			WHEN 'name' THEN CAST(u.name AS char)
			WHEN 'surname' THEN CAST(u.surname AS char)
			WHEN 'patronymic' THEN CAST(u.patronymic AS char)
			WHEN 'username' THEN CAST(u.username AS char)
			WHEN 'birthday' THEN CAST(u.birthday AS char)
			WHEN 'status' THEN CAST(u.status AS char)
			ELSE CAST(u.created_at AS char)
		END
	END DESC
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
	r.service_id = sqlc.arg (service_id)
	AND (
		CAST(sqlc.narg (search) AS char) IS NULL
		OR u.name LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
		OR u.surname LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
		OR u.patronymic LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
		OR u.username LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	);

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
