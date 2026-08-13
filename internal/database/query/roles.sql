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
	roles r
WHERE
	(
		CAST(sqlc.narg (search) AS char) IS NULL
		OR r.name LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
		OR r.description LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	)
	AND (
		sqlc.narg (is_global) IS NULL
		OR r.is_global = sqlc.narg (is_global)
	);

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
	roles r
WHERE
	r.service_id = sqlc.arg (service_id)
	AND (
		CAST(sqlc.narg (search) AS char) IS NULL
		OR r.name LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
		OR r.description LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	);

-- name: ListRolesWithCounts :many
-- search: поиск по name/description. is_global: true/false — фильтр по
-- колонке is_global (NULL/не передан — без фильтра). sort_by: name/description
-- /created_at (иначе created_at), sort_dir: asc (иначе desc) — см. query.QuerySort.
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
	(
		CAST(sqlc.narg (search) AS char) IS NULL
		OR r.name LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
		OR r.description LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	)
	AND (
		sqlc.narg (is_global) IS NULL
		OR r.is_global = sqlc.narg (is_global)
	)
ORDER BY
	CASE
		WHEN CAST(sqlc.arg (sort_dir) AS char) = 'asc' THEN CASE CAST(sqlc.arg (sort_by) AS char)
			WHEN 'name' THEN CAST(r.name AS char)
			WHEN 'description' THEN CAST(r.description AS char)
			ELSE CAST(r.created_at AS char)
		END
	END ASC,
	CASE
		WHEN CAST(sqlc.arg (sort_dir) AS char) = 'desc' THEN CASE CAST(sqlc.arg (sort_by) AS char)
			WHEN 'name' THEN CAST(r.name AS char)
			WHEN 'description' THEN CAST(r.description AS char)
			ELSE CAST(r.created_at AS char)
		END
	END DESC
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
	r.service_id = sqlc.arg (service_id)
	AND (
		CAST(sqlc.narg (search) AS char) IS NULL
		OR r.name LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
		OR r.description LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	)
ORDER BY
	CASE
		WHEN CAST(sqlc.arg (sort_dir) AS char) = 'asc' THEN CASE CAST(sqlc.arg (sort_by) AS char)
			WHEN 'name' THEN CAST(r.name AS char)
			WHEN 'description' THEN CAST(r.description AS char)
			ELSE CAST(r.created_at AS char)
		END
	END ASC,
	CASE
		WHEN CAST(sqlc.arg (sort_dir) AS char) = 'desc' THEN CASE CAST(sqlc.arg (sort_by) AS char)
			WHEN 'name' THEN CAST(r.name AS char)
			WHEN 'description' THEN CAST(r.description AS char)
			ELSE CAST(r.created_at AS char)
		END
	END DESC
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
