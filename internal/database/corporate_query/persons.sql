-- name: CreatePerson :exec
INSERT INTO
	persons (id, first_name, last_name, patronymic, birthday, snils, phone, email, auth_user_id)
VALUES
	(?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetPersonByID :one
SELECT
	id,
	first_name,
	last_name,
	patronymic,
	birthday,
	snils,
	phone,
	email,
	auth_user_id,
	created_at
FROM
	persons
WHERE
	id = ?;

-- name: UpdatePersonOnLink :exec
-- Используется в invite.Service.Accept, когда инвайт был выпущен с уже
-- существующим person_id: линкует новый auth-аккаунт к существующему
-- профилю (вместо создания нового) и одновременно перезаписывает
-- редактируемые поля тем, что пользователь подтвердил/поправил в форме
-- (см. ValidateResponse.Person — оттуда форма была предзаполнена).
UPDATE persons
SET
	first_name = ?,
	last_name = ?,
	patronymic = ?,
	birthday = ?,
	phone = ?,
	email = ?,
	auth_user_id = ?
WHERE
	id = ?;

-- name: GetPersonByAuthUserID :one
SELECT
	id,
	first_name,
	last_name,
	patronymic,
	birthday,
	snils,
	phone,
	email,
	auth_user_id,
	created_at
FROM
	persons
WHERE
	auth_user_id = ?;
