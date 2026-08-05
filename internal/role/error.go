package role

import "errors"

var (
	ErrNameRequired        = errors.New("name обязателен")
	ErrDescriptionRequired = errors.New("description обязателен")

	ErrNotFound = errors.New("роль не найдена")

	ErrPermissionNotFound = errors.New("разрешение не найдено")

	ErrPermissionAlreadyAssigned = errors.New("роль уже имеет это разрешение")

	ErrPermissionNotAssigned = errors.New("роль не имеет это разрешение")
)
