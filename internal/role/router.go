package role

import (
	"permisson/internal/pkg/apidoc"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes — аналог roles_router = APIRouter(prefix="/roles", tags=["Roles"]).
//
// Как и в остальных сущностях, require_permission пока не подключён — роуты
// работают, но без проверки прав до появления middleware.RequirePermission.
func RegisterRoutes(router fiber.Router, schema *ftonic.Adapter, h *Handler) {
	group := router.Group("/roles")

	ftonic.For[apidoc.Pagination, ListResponse](schema).
		GET(group, "/", h.List, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список ролей",
			Description: "Требует roles:read.",
			Tags:        []string{"Roles"},
		}))

	ftonic.For[ListByServiceRequest, ListResponse](schema).
		GET(group, "/service/:service_id", h.ListByServiceID, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список ролей по service_id",
			Description: "Требует roles:read.",
			Tags:        []string{"Roles"},
		}))

	ftonic.For[AddPermissionRequest, RoleResponse](schema).
		POST(group, "/perm/add", h.AddPermission, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Добавить разрешение для роли",
			Description: "Требует roles.perm:edit.",
			Tags:        []string{"Roles"},
		}))

	ftonic.For[AddPermissionRequest, RoleResponse](schema).
		POST(group, "/perm/remove", h.RemovePermission, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Удалить разрешение у роли",
			Description: "Требует roles.perm:edit.",
			Tags:        []string{"Roles"},
		}))

	ftonic.For[UpsertRequest, RoleResponse](schema).
		POST(group, "/create", h.Create, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Создать роль",
			Description: "Требует roles:create.",
			Tags:        []string{"Roles"},
		}))

	ftonic.For[UpsertRequest, RoleResponse](schema).
		PUT(group, "/:role_id", h.Update, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Редактировать роль",
			Description: "Требует roles:update.",
			Tags:        []string{"Roles"},
		}))

	ftonic.For[apidoc.Empty, DetailedResponse](schema).
		GET(group, "/:role_id/detailed", h.Detailed, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Подробная информация о роли",
			Description: "Разрешения, сгруппированные по сервисам, + список пользователей с этой ролью. Требует roles:read.",
			Tags:        []string{"Roles"},
		}))

	ftonic.For[apidoc.Empty, DeleteResponse](schema).
		DELETE(group, "/:role_id", h.Delete, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Удалить роль",
			Description: "Требует roles:delete.",
			Tags:        []string{"Roles"},
		}))
}
