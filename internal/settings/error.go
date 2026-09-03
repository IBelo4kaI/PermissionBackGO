package settings

import "errors"

var (
	ErrHostRequired        = errors.New("host обязателен")
	ErrPortRequired        = errors.New("port обязателен")
	ErrUsernameRequired    = errors.New("username обязателен")
	ErrFromAddressRequired = errors.New("from_address обязателен")
	ErrPasswordRequired    = errors.New("password обязателен при первоначальной настройке")

	// ErrNotConfigured — используется только GetForSending: отправить письмо
	// нечем, пока админ не задал SMTP хотя бы раз.
	ErrNotConfigured = errors.New("SMTP не настроен")
)
