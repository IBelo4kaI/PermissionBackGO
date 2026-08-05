package role

import "errors"

var (
	ErrNameRequired        = errors.New("name обязателен")
	ErrDescriptionRequired = errors.New("description обязателен")

	// ErrNotFound — аналог HTTPException(404, "Роль не найдена").
	ErrNotFound = errors.New("роль не найдена")

	// ErrPermissionNotFound — аналог HTTPException(404, "Разрешение не найдено").
	// В Python-версии RoleService.permission_add при отсутствии разрешения
	// ошибочно возвращал текст "Роль не найдена" (copy-paste баг из строки выше) —
	// здесь текст исправлен на корректный.
	ErrPermissionNotFound = errors.New("разрешение не найдено")

	// ErrPermissionAlreadyAssigned — аналог HTTPException(400, "Роль уже имеет это разрешение").
	ErrPermissionAlreadyAssigned = errors.New("роль уже имеет это разрешение")

	// ErrPermissionNotAssigned — аналог HTTPException(400, "Роль не имеет это разрешение").
	ErrPermissionNotAssigned = errors.New("роль не имеет это разрешение")
)
