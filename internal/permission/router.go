package permission

import (
	"permisson/internal/pkg/apidoc"
	"permisson/internal/pkg/response"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, h *Handler, schema *ftonic.Adapter) {
	group := router.Group("/permissions")

	ftonic.For[apidoc.Pagination, response.Page[Permission]](schema).
		GET(group, "/", h.List, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список всех разрешений",
			Description: "Требует perm:read (проверка прав будет добавлена вместе с middleware).",
			Tags:        []string{"Разрешения"},
		}))

	ftonic.For[ListByServiceRequest, []Permission](schema).
		GET(group, "/:service_id", h.ListByServiceID, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список разрешений для сервиса",
			Description: "Требует perm:read.",
			Tags:        []string{"Разрешения"},
		}))

	ftonic.For[UpsertRequest, Permission](schema).
		POST(group, "/create", h.Create, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Создать разрешение",
			Description: "Требует perm:create.",
			Tags:        []string{"Разрешения"},
		}))

	ftonic.For[UpdatePermissionRequest, Permission](schema).
		PUT(group, "/:permission_id", h.Update, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Редактировать разрешение",
			Description: "Требует perm:update.",
			Tags:        []string{"Разрешения"},
		}))

	ftonic.For[PermissionIDRequest, DeleteResponse](schema).
		DELETE(group, "/:permission_id", h.Delete, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Удалить разрешение",
			Description: "Требует perm:delete.",
			Tags:        []string{"Разрешения"},
		}))
}
