package permission

import "errors"

var (
	ErrCodeRequired        = errors.New("code обязателен")
	ErrNameRequired        = errors.New("name обязателен")
	ErrDescriptionRequired = errors.New("description обязателен")

	ErrCodeExists = errors.New("разрешение с таким кодом уже существует")

	ErrNotFound = errors.New("разрешение не найдено")
)
