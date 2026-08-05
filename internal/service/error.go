package service

import "errors"

var (
	ErrNameRequired        = errors.New("name обязателен")
	ErrNameTooLong         = errors.New("name не может быть длиннее 100 символов")
	ErrDescriptionRequired = errors.New("description обязателен")
	ErrPrefixRequired      = errors.New("prefix обязателен")
	ErrPrefixTooLong       = errors.New("prefix не может быть длиннее 5 символов")

	// ErrNotFound — аналог HTTPException(404, "Сервис не найден").
	ErrNotFound = errors.New("сервис не найден")
)
