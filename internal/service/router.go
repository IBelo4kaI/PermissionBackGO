package service

import (
	"permisson/internal/pkg/apidoc"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes — аналог service_router = APIRouter(prefix="/services", tags=["Services"]).
// Как и в permission-сущности, require_permission пока не подключён — роуты
// работают, но без проверки прав до появления middleware.RequirePermission.
func RegisterRoutes(router fiber.Router, h *Handler, schema *ftonic.Adapter) {
	group := router.Group("/services")

	ftonic.For[apidoc.Pagination, ListResponse](schema).
		GET(group, "/", h.List, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список сервисов",
			Description: "Требует services:read_all.",
			Tags:        []string{"Services"},
		}))

	ftonic.For[UpsertRequest, ServiceResponse](schema).
		POST(group, "/create", h.Create, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Создать сервис",
			Description: "Требует services:create.",
			Tags:        []string{"Services"},
		}))

	ftonic.For[UpsertRequest, ServiceResponse](schema).
		PUT(group, "/:service_id", h.Update, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Редактировать сервис",
			Description: "Требует services:update.",
			Tags:        []string{"Services"},
		}))

	ftonic.For[apidoc.Empty, []AccessResponse](schema).
		GET(group, "/user-accessible", h.ListUserAccessible, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Сервисы, доступные текущему пользователю",
			Description: "Требует валидную пользовательскую cookie-сессию (не API-ключ).",
			Tags:        []string{"Services"},
		}))

	ftonic.For[apidoc.Empty, APIKeyResponse](schema).
		POST(group, "/:service_id/api-key", h.IssueAPIKey, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Выпустить (или перевыпустить) API-ключ сервиса",
			Description: "Требует services:update. Сырой ключ возвращается один раз.",
			Tags:        []string{"Services"},
		}))

	ftonic.For[apidoc.Empty, apidoc.Empty](schema).
		DELETE(group, "/:service_id/api-key", h.RevokeAPIKey, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Отозвать API-ключ сервиса",
			Description: "Требует services:update.",
			Tags:        []string{"Services"},
		}))
}
