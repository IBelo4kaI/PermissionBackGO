-- name: CreateEmployment :exec
INSERT INTO
	employments (id, person_id, company_id, department_id, position_id, status, hired_at)
VALUES
	(?, ?, ?, ?, ?, ?, ?);

-- name: GetEmploymentByPersonID :many
SELECT
	id,
	person_id,
	company_id,
	department_id,
	position_id,
	phone_work,
	email_work,
	status,
	use_timetrack,
	hired_at,
	dismissed_at,
	role_department,
	created_at
FROM
	employments
WHERE
	person_id = ?;
