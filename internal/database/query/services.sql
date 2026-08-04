-- name: GetServiceByID :one
SELECT id, name, description, image_url, url, theme, prefix, api_key_hash, created_at
FROM services
WHERE id = ?;

-- name: GetServiceByAPIKeyHash :one
SELECT id, name, description, image_url, url, theme, prefix, api_key_hash, created_at
FROM services
WHERE api_key_hash = ?;

-- name: ListServices :many
SELECT id, name, description, image_url, url, theme, prefix, api_key_hash, created_at
FROM services
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountServices :one
SELECT COUNT(*) AS total
FROM services;

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
