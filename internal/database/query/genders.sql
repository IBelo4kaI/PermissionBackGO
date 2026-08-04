-- name: GetGenderByID :one
SELECT id, name
FROM genders
WHERE id = ?;

-- name: ListGenders :many
SELECT id, name
FROM genders
ORDER BY id;
