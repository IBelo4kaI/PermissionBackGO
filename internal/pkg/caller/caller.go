// Package caller описывает единую идентичность запроса: пользователь по
// сессии или сервис по API-ключу. Вынесен в отдельный пакет, чтобы его могли
// использовать и auth (определение), и middleware (потребление) без
// циклической зависимости auth ↔ middleware.
package caller

type Type string

const (
	User    Type = "user"
	Service Type = "service"
)

type Caller struct {
	Type      Type
	UserID    string
	ServiceID string
}
