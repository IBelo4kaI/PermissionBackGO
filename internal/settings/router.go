package settings

import (
	middlewares "permisson/internal/middleware"
	"permisson/internal/pkg/apidoc"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, h *Handler, schema *ftonic.Adapter, require middlewares.Require) {
	group := router.Group("/settings")

	ftonic.For[apidoc.Empty, SMTPSettingsResponse](schema).
		GET(group, "/smtp", require("settings", "read"), h.GetSMTP, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить настройки SMTP",
			Description: "Требует settings:read. Пароль в ответе никогда не возвращается.",
			Tags:        []string{"Настройки"},
		}))

	ftonic.For[UpsertSMTPSettingsRequest, SMTPSettingsResponse](schema).
		PUT(group, "/smtp", require("settings", "update"), h.UpsertSMTP, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Задать настройки SMTP",
			Description: "Требует settings:update. Password опционален при обновлении — не передан, прежний пароль сохраняется.",
			Tags:        []string{"Настройки"},
		}))
}
