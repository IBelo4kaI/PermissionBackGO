-- name: GetPermissionByID :one
SELECT
	p.id,
	p.service_id,
	p.code,
	p.name,
	p.description,
	p.created_at,
	s.name AS service_name
FROM
	permissions p
	LEFT JOIN services s ON s.id = p.service_id
WHERE
	p.id = ?;

-- name: GetPermissionByCode :one
SELECT
	p.id,
	p.service_id,
	p.code,
	p.name,
	p.description,
	p.created_at,
	s.name AS service_name
FROM
	permissions p
	LEFT JOIN services s ON s.id = p.service_id
WHERE
	p.code = ?;

-- name: ListPermissions :many
SELECT
	p.id,
	p.service_id,
	p.code,
	p.name,
	p.description,
	p.created_at,
	s.name AS service_name
FROM
	permissions p
	LEFT JOIN services s ON s.id = p.service_id
ORDER BY
	p.created_at DESC
LIMIT
	?
OFFSET
	?;

-- name: CountPermissions :one
SELECT
	COUNT(*) AS total
FROM
	permissions;

-- name: ListPermissionsByServiceID :many
SELECT
	p.id,
	p.service_id,
	p.code,
	p.name,
	p.description,
	p.created_at,
	s.name AS service_name
FROM
	permissions p
	LEFT JOIN services s ON s.id = p.service_id
WHERE
	p.service_id = ?
ORDER BY
	p.created_at DESC
LIMIT
	?
OFFSET
	?;

-- name: CountPermissionsByServiceID :one
SELECT
	COUNT(*) AS total
FROM
	permissions
WHERE
	service_id = ?;

-- name: ListPermissionsByUserID :many
SELECT DISTINCT
	p.id,
	p.service_id,
	p.code,
	p.name,
	p.description,
	p.created_at
FROM
	permissions p
	JOIN role_permissions rp ON rp.permission_id = p.id
	JOIN user_roles ur ON ur.role_id = rp.role_id
WHERE
	ur.user_id = ?;

-- name: ListPermissionsByUserIDAndServiceID :many
-- Разрешения пользователя для конкретного сервиса с учётом wildcard-кодов
-- вида "all:all:all" и "<service_name>:all:all". sqlc.arg + CAST(...AS CHAR),
-- чтобы sqlc корректно типизировал повторяющийся параметр service_name.
SELECT DISTINCT
	p.id,
	p.service_id,
	p.code,
	p.name,
	p.description,
	p.created_at
FROM
	permissions p
	JOIN role_permissions rp ON rp.permission_id = p.id
	JOIN user_roles ur ON ur.role_id = rp.role_id
WHERE
	ur.user_id = sqlc.arg (user_id)
	AND (
		p.service_id = sqlc.arg (service_id)
		OR p.code like CONCAT(CAST(sqlc.arg (service_name) AS char), ':%')
		OR p.code like 'all:%'
	);

-- name: ExistsUserPermission :one
-- Проверка наличия разрешения у пользователя с учётом всех wildcard-комбинаций
-- сегментов "service:entity:action". Каждый sqlc.arg(...) используется несколько
-- раз, но CAST(... AS CHAR) даёт sqlc однозначный тип (string), поэтому в
-- сгенерированных параметрах будет ровно 4 понятных поля, а не CONCAT_2..CONCAT_12.
SELECT
	EXISTS (
		SELECT
			1
		FROM
			permissions p
			JOIN role_permissions rp ON rp.permission_id = p.id
			JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE
			ur.user_id = sqlc.arg (user_id)
			AND p.code IN (
				CONCAT(CAST(sqlc.arg (service) AS char), ':', CAST(sqlc.arg (entity) AS char), ':', CAST(sqlc.arg (action) AS char)),
				CONCAT(CAST(sqlc.arg (service) AS char), ':', CAST(sqlc.arg (entity) AS char), ':all'),
				CONCAT(CAST(sqlc.arg (service) AS char), ':all:', CAST(sqlc.arg (action) AS char)),
				CONCAT(CAST(sqlc.arg (service) AS char), ':all:all'),
				CONCAT('all:', CAST(sqlc.arg (entity) AS char), ':', CAST(sqlc.arg (action) AS char)),
				CONCAT('all:', CAST(sqlc.arg (entity) AS char), ':all'),
				CONCAT('all:all:', CAST(sqlc.arg (action) AS char)),
				'all:all:all'
			)
		LIMIT
			1
	) AS has_permission;

-- name: ListAllPermissions :many
SELECT
	p.id,
	p.service_id,
	p.code,
	p.name,
	p.description,
	p.created_at,
	s.name AS service_name
FROM
	permissions p
	LEFT JOIN services s ON s.id = p.service_id;

-- name: ListAllPermissionsByServiceID :many
SELECT
	p.id,
	p.service_id,
	p.code,
	p.name,
	p.description,
	p.created_at,
	s.name AS service_name
FROM
	permissions p
	LEFT JOIN services s ON s.id = p.service_id
WHERE
	p.service_id = ?;

-- name: CreatePermission :exec
INSERT INTO
	permissions (id, service_id, code, name, description, created_at)
VALUES
	(?, ?, ?, ?, ?, NOW());

-- name: UpdatePermission :exec
UPDATE permissions
SET
	service_id = ?,
	code = ?,
	name = ?,
	description = ?
WHERE
	id = ?;

-- name: DeletePermission :exec
DELETE FROM permissions
WHERE
	id = ?;
