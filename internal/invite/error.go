package invite

import "errors"

var (
	ErrNotFound     = errors.New("инвайт не найден")
	ErrCodeNotFound = errors.New("инвайт не найден или недействителен")

	ErrInviteUsed    = errors.New("инвайт уже использован")
	ErrInviteRevoked = errors.New("инвайт отозван")
	ErrInviteExpired = errors.New("срок действия инвайта истёк")
)
