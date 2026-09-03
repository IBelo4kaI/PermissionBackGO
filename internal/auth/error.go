package auth

import "errors"

var (
	ErrEmailTaken          = errors.New("Почта уже используется")
	ErrInvalidCredentials  = errors.New("Неверная почта или пароль")
	ErrUserInactive        = errors.New("Аккаунт пользователя неактивен, обратитесь к администратору")
	ErrUsernameRequired    = errors.New("Почта обязательна")
	ErrPasswordRequired    = errors.New("Пароль обязателен")
	ErrInvalidSession      = errors.New("Недействительная сессия")
	ErrSessionTokenMissing = errors.New("Session token missing")

	ErrTokenRequired       = errors.New("token обязателен")
	ErrNewPasswordRequired = errors.New("new_password обязателен")
	ErrInvalidResetToken   = errors.New("Ссылка недействительна, истекла или уже использована")
)
