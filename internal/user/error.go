package user

import "errors"

var (
	ErrNameRequired     = errors.New("имя обязательно")
	ErrSurnameRequired  = errors.New("фамилия обязательна")
	ErrUsernameRequired = errors.New("username обязателен")
	ErrPasswordRequired = errors.New("пароль обязателен")
	ErrBirthdayRequired = errors.New("дата рождения обязательна")
	ErrGenderRequired   = errors.New("gender_id обязателен")

	ErrNotFound = errors.New("пользователь не найден")

	ErrGenderNotFound = errors.New("не существующий идентификатор Gender")

	ErrUsernameExists = errors.New("пользователь с таким именем уже существует")

	ErrRoleNotFound = errors.New("роль не найдена")

	ErrRoleAlreadyAssigned = errors.New("пользователь уже имеет эту роль")

	ErrRoleNotAssigned = errors.New("пользователь не имеет этой роли")
)
