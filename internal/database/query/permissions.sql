-- name: GetPermissionByID :one
SELECT id, service_id, code, name, description, created_at
FROM permissions
WHERE id = ?;

-- name: GetPermissionByCode :one
SELECT id, service_id, code, name, description, created_at
FROM permissions
WHERE code = ?;

-- name: ListPermissions :many
SELECT id, service_id, code, name, description, created_at
FROM permissions
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountPermissions :one
SELECT COUNT(*) AS total
FROM permissions;

-- name: ListPermissionsByServiceID :many
SELECT id, service_id, code, name, description, created_at
FROM permissions
WHERE service_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: CountPermissionsByServiceID :one
SELECT COUNT(*) AS total
FROM permissions
WHERE service_id = ?;

-- name: ListPermissionsByUserID :many
-- Все разрешения, доступные пользователю через все его роли.
SELECT DISTINCT p.id, p.service_id, p.code, p.name, p.description, p.created_at
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
JOIN user_roles ur ON ur.role_id = rp.role_id
WHERE ur.user_id = ?;

-- name: ListPermissionsByUserIDAndServiceID :many
-- Разрешения пользователя для конкретного сервиса с учётом wildcard-кодов
-- вида "all:all:all" и "<service_name>:all:all".
SELECT DISTINCT p.id, p.service_id, p.code, p.name, p.description, p.created_at
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
JOIN user_roles ur ON ur.role_id = rp.role_id
WHERE ur.user_id = ?
  AND (
    p.service_id = ?
    OR p.code LIKE CONCAT(?, ':%')
    OR p.code LIKE 'all:%'
  );

-- name: ExistsUserPermission :one
-- Проверка наличия разрешения у пользователя с учётом всех wildcard-комбинаций
-- сегментов "service:entity:action".
SELECT EXISTS (
    SELECT 1
    FROM permissions p
    JOIN role_permissions rp ON rp.permission_id = p.id
    JOIN user_roles ur ON ur.role_id = rp.role_id
    WHERE ur.user_id = ?
      AND p.code IN (
        CONCAT(?, ':', ?, ':', ?),
        CONCAT(?, ':', ?, ':all'),
        CONCAT(?, ':all:', ?),
        CONCAT(?, ':all:all'),
        CONCAT('all:', ?, ':', ?),
        CONCAT('all:', ?, ':all'),
        CONCAT('all:all:', ?),
        'all:all:all'
      )
    LIMIT 1
) AS has_permission;

-- name: ListAllPermissions :many
-- Без пагинации: используется, например, при сборе информации о глобальной роли.
SELECT p.id, p.service_id, p.code, p.name, p.description, p.created_at, s.name AS service_name
FROM permissions p
LEFT JOIN services s ON s.id = p.service_id;

-- name: ListAllPermissionsByServiceID :many
-- Без пагинации: используется при сборе информации о роли, привязанной к сервису.
SELECT p.id, p.service_id, p.code, p.name, p.description, p.created_at, s.name AS service_name
FROM permissions p
LEFT JOIN services s ON s.id = p.service_id
WHERE p.service_id = ?;

-- name: CreatePermission :exec
INSERT INTO permissions (id, service_id, code, name, description, created_at)
VALUES (?, ?, ?, ?, ?, NOW());

-- name: UpdatePermission :exec
UPDATE permissions
SET service_id  = ?,
    code        = ?,
    name        = ?,
    description = ?
WHERE id = ?;

-- name: DeletePermission :exec
DELETE FROM permissions
WHERE id = ?;
