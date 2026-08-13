package service

import (
	middlewares "permisson/internal/middleware"
	"permisson/internal/pkg/apidoc"
	"permisson/internal/pkg/response"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, h *Handler, schema *ftonic.Adapter, require middlewares.Require) {
	group := router.Group("/services")

	ftonic.For[apidoc.Pagination, response.Page[ServiceResponse]](schema).
		GET(group, "/", require("services", "read_all"), h.List, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список сервисов",
			Description: "Требует services:read_all.",
			Tags:        []string{"Сервисы"},
		}))

	ftonic.For[UpsertRequest, ServiceResponse](schema).
		POST(group, "/create", require("services", "create"), h.Create, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Создать сервис",
			Description: "Требует services:create.",
			Tags:        []string{"Сервисы"},
		}))

	ftonic.For[UpdateServiceRequest, ServiceResponse](schema).
		PUT(group, "/:service_id", require("services", "update"), h.Update, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Редактировать сервис",
			Description: "Требует services:update.",
			Tags:        []string{"Сервисы"},
		}))

	ftonic.For[ListUserAccessibleRequest, []AccessResponse](schema).
		GET(group, "/user-accessible", h.ListUserAccessible, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Сервисы, доступные текущему пользователю",
			Description: "Требует валидную пользовательскую cookie-сессию (не API-ключ).",
			Tags:        []string{"Сервисы"},
		}))

	ftonic.For[ServiceIDRequest, APIKeyResponse](schema).
		POST(group, "/:service_id/api-key", h.IssueAPIKey, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Выпустить (или перевыпустить) API-ключ сервиса",
			Description: "Требует services:update. Сырой ключ возвращается один раз.",
			Tags:        []string{"Сервисы"},
		}))

	ftonic.For[ServiceIDRequest, apidoc.Empty](schema).
		DELETE(group, "/:service_id/api-key", require("services", "update"), h.RevokeAPIKey, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Отозвать API-ключ сервиса",
			Description: "Требует services:update.",
			Tags:        []string{"Сервисы"},
		}))
}
