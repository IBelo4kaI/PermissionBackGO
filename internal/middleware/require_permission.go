package middlewares

import (
	"context"

	"github.com/gofiber/fiber/v3"

	"permisson/internal/pkg/caller"
)

// SessionResolver определяет идентичность запроса по cookie "session".
// Реализуется *auth.Service; интерфейс нужен, чтобы не создавать циклическую
// зависимость middleware ↔ auth (auth/router.go использует middlewares.Bind).
type SessionResolver interface {
	CallerFrom(ctx context.Context, rawToken string) (caller.Caller, error)
}

// PermissionChecker — проверка наличия разрешения у пользователя.
// Реализуется *permission.Service; интерфейс нужен, чтобы не создавать
// циклическую зависимость middleware ↔ permission.
type PermissionChecker interface {
	ExistsForUser(ctx context.Context, userID, service, entity, action string) (bool, error)
}

// Require — фабрика middleware проверки прав для конкретного роута.
// Роутеры объявляют права декларативно: require("users", "read").
type Require func(entity, action string) fiber.Handler

// NewRequire связывает аутентификацию и проверку прав, чтобы не тащить
// сервисы в каждый роутер отдельно.
func NewRequire(authService SessionResolver, permService PermissionChecker) Require {
	return func(entity, action string) fiber.Handler {
		return RequirePermission(authService, permService, entity, action)
	}
}

// RequirePermission — аналог Python require_permission(entity, action):
// сервис с валидным API-ключом полностью доверенный (гранулярные права не
// проверяем), для пользователя проверяется наличие разрешения
// "perm:<entity>:<action>" с учётом всех wildcard-комбинаций.
func RequirePermission(authService SessionResolver, permService PermissionChecker, entity, action string) fiber.Handler {
	return func(c fiber.Ctx) error {
		callerInfo, err := authService.CallerFrom(c.Context(), c.Cookies("session"))
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		if callerInfo.Type == caller.Service {
			return c.Next()
		}
		if callerInfo.UserID == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		ok, err := permService.ExistsForUser(c.Context(), callerInfo.UserID, "perm", entity, action)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
		}
		if !ok {
			return fiber.NewError(fiber.StatusForbidden, "Нет доступа")
		}
		return c.Next()
	}
}
