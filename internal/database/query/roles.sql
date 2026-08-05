-- name: GetRoleByID :one
SELECT
	id,
	service_id,
	name,
	description,
	is_global,
	created_at
FROM
	roles
WHERE
	id = ?;

-- name: ListRoles :many
SELECT
	id,
	service_id,
	name,
	description,
	is_global,
	created_at
FROM
	roles
ORDER BY
	created_at DESC
LIMIT
	?
OFFSET
	?;

-- name: CountRoles :one
SELECT
	COUNT(*) AS total
FROM
	roles;

-- name: ListRolesByServiceID :many
SELECT
	id,
	service_id,
	name,
	description,
	is_global,
	created_at
FROM
	roles
WHERE
	service_id = ?
ORDER BY
	created_at DESC
LIMIT
	?
OFFSET
	?;

-- name: CountRolesByServiceID :one
SELECT
	COUNT(*) AS total
FROM
	roles
WHERE
	service_id = ?;

-- name: ListRolesWithCounts :many
SELECT
	r.id,
	r.service_id,
	r.name,
	r.description,
	r.is_global,
	r.created_at,
	COALESCE(uc.user_count, 0) AS user_count,
	COALESCE(pc.permission_count, 0) AS permission_count
FROM
	roles r
	LEFT JOIN (
		SELECT
			role_id,
			COUNT(user_id) AS user_count
		FROM
			user_roles
		GROUP BY
			role_id
	) uc ON uc.role_id = r.id
	LEFT JOIN (
		SELECT
			role_id,
			COUNT(permission_id) AS permission_count
		FROM
			role_permissions
		GROUP BY
			role_id
	) pc ON pc.role_id = r.id
ORDER BY
	r.created_at DESC
LIMIT
	?
OFFSET
	?;

-- name: ListRolesWithCountsByServiceID :many
SELECT
	r.id,
	r.service_id,
	r.name,
	r.description,
	r.is_global,
	r.created_at,
	COALESCE(uc.user_count, 0) AS user_count,
	COALESCE(pc.permission_count, 0) AS permission_count
FROM
	roles r
	LEFT JOIN (
		SELECT
			role_id,
			COUNT(user_id) AS user_count
		FROM
			user_roles
		GROUP BY
			role_id
	) uc ON uc.role_id = r.id
	LEFT JOIN (
		SELECT
			role_id,
			COUNT(permission_id) AS permission_count
		FROM
			role_permissions
		GROUP BY
			role_id
	) pc ON pc.role_id = r.id
WHERE
	r.service_id = ?
ORDER BY
	r.created_at DESC
LIMIT
	?
OFFSET
	?;

-- name: CreateRole :exec
INSERT INTO
	roles (id, service_id, name, description, is_global, created_at)
VALUES
	(?, ?, ?, ?, ?, NOW());

-- name: UpdateRole :exec
UPDATE roles
SET
	service_id = ?,
	name = ?,
	description = ?,
	is_global = ?
WHERE
	id = ?;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE
	id = ?;

-- name: AddPermissionToRole :exec
INSERT INTO
	role_permissions (role_id, permission_id, granted_at)
VALUES
	(?, ?, NOW());

-- name: RemovePermissionFromRole :exec
DELETE FROM role_permissions
WHERE
	role_id = ?
	AND permission_id = ?;

-- name: ListPermissionsForRole :many
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
	JOIN role_permissions rp ON rp.permission_id = p.id
	LEFT JOIN services s ON s.id = p.service_id
WHERE
	rp.role_id = ?;

-- name: ListUsedPermissionIDsForRole :many
SELECT
	permission_id
FROM
	role_permissions
WHERE
	role_id = ?;

-- name: ListUsersWithRole :many
SELECT
	u.id,
	u.name,
	u.surname,
	u.patronymic,
	u.username,
	u.gender_id,
	u.birthday,
	u.status,
	u.created_at
FROM
	users u
	JOIN user_roles ur ON ur.user_id = u.id
WHERE
	ur.role_id = ?;
