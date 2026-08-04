package auth

import "errors"

var (
	ErrEmailTaken         = errors.New("Почта уже используется")
	ErrInvalidCredentials = errors.New("Неверная почта или пароль")
	ErrUserInactive       = errors.New("Пользователь не активный")
	ErrUsernameRequired   = errors.New("Почта обязательна")
	ErrPasswordRequired   = errors.New("Пароль обязателен")
)
