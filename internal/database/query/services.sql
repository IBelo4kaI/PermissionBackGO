-- name: GetServiceByID :one
SELECT id, name, description, image_url, url, theme, prefix, api_key_hash, created_at
FROM services
WHERE id = ?;

-- name: GetServiceByAPIKeyHash :one
SELECT id, name, description, image_url, url, theme, prefix, api_key_hash, created_at
FROM services
WHERE api_key_hash = ?;

-- name: ListServices :many
-- search: поиск по name/description/prefix. sort_by: name/description/prefix/
-- created_at (иначе created_at), sort_dir: asc (иначе desc) — см. query.QuerySort.
SELECT id, name, description, image_url, url, theme, prefix, api_key_hash, created_at
FROM services s
WHERE
	CAST(sqlc.narg (search) AS char) IS NULL
	OR s.name LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	OR s.description LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	OR s.prefix LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
ORDER BY
	CASE
		WHEN CAST(sqlc.arg (sort_dir) AS char) = 'asc' THEN CASE CAST(sqlc.arg (sort_by) AS char)
			WHEN 'name' THEN CAST(s.name AS char)
			WHEN 'description' THEN CAST(s.description AS char)
			WHEN 'prefix' THEN CAST(s.prefix AS char)
			ELSE CAST(s.created_at AS char)
		END
	END ASC,
	CASE
		WHEN CAST(sqlc.arg (sort_dir) AS char) = 'desc' THEN CASE CAST(sqlc.arg (sort_by) AS char)
			WHEN 'name' THEN CAST(s.name AS char)
			WHEN 'description' THEN CAST(s.description AS char)
			WHEN 'prefix' THEN CAST(s.prefix AS char)
			ELSE CAST(s.created_at AS char)
		END
	END DESC
LIMIT ? OFFSET ?;

-- name: CountServices :one
SELECT COUNT(*) AS total
FROM services s
WHERE
	CAST(sqlc.narg (search) AS char) IS NULL
	OR s.name LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	OR s.description LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%')
	OR s.prefix LIKE CONCAT('%', CAST(sqlc.narg (search) AS char), '%');

-- name: CreateService :exec
INSERT INTO services (id, name, description, image_url, url, theme, prefix, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, NOW());

-- name: UpdateService :exec
UPDATE services
SET name        = ?,
    description = ?,
    image_url   = ?,
    url         = ?,
    theme       = ?,
    prefix      = ?
WHERE id = ?;

-- name: SetServiceAPIKeyHash :exec
UPDATE services
SET api_key_hash = ?
WHERE id = ?;

-- name: RevokeServiceAPIKey :exec
UPDATE services
SET api_key_hash = NULL
WHERE id = ?;

-- name: DeleteService :exec
DELETE FROM services
WHERE id = ?;

-- name: ListServicesByUserID :many
-- Сервисы, доступные пользователю через его роли:
-- либо роль привязана к конкретному сервису, либо роль глобальная (service_id IS NULL).
SELECT DISTINCT s.id, s.name, s.description, s.image_url, s.url, s.theme, s.prefix, s.api_key_hash, s.created_at
FROM services s
JOIN roles r ON (r.service_id = s.id OR r.service_id IS NULL)
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = ?;
