package auth

import (
	middlewares "permisson/internal/middleware"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, h *Handler, schema *ftonic.Adapter) {
	group := router.Group("/auth")

	ftonic.For[LoginRequest, LoginResult](schema).
		POST(group, "/login", middlewares.Bind[LoginRequest], h.Login, ftonic.WithOperation(docs.OperationObject{
			Summary: "Авторизация в системе",
			Tags:    []string{"Авторизация"},
		}))

	ftonic.For[any, ValidateSessionResponse](schema).
		POST(group, "/validate-session", h.ValidateSession, ftonic.WithOperation(docs.OperationObject{
			Summary: "Проверка токена на валидность",
			Tags:    []string{"Авторизация"},
		}))

	ftonic.For[any, LogoutResponse](schema).
		POST(group, "/logout", h.Logout, ftonic.WithOperation(docs.OperationObject{
			Summary: "Выход из системы и удаление сессии",
			Tags:    []string{"Авторизация"},
		}))

	ftonic.For[ForgotPasswordRequest, ForgotPasswordResponse](schema).
		POST(group, "/forgot-password", h.ForgotPassword, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Запросить восстановление пароля",
			Description: "Публичный эндпоинт. Ответ одинаковый независимо от того, существует ли пользователь.",
			Tags:        []string{"Авторизация"},
		}))

	ftonic.For[ResetPasswordRequest, ResetPasswordResponse](schema).
		POST(group, "/reset-password", h.ResetPassword, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Установить новый пароль по ссылке из письма",
			Description: "Публичный эндпоинт. Удаляет все сессии пользователя после смены пароля.",
			Tags:        []string{"Авторизация"},
		}))
}
