package role

import (
	middlewares "permisson/internal/middleware"
	"permisson/internal/pkg/apidoc"
	"permisson/internal/pkg/response"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, h *Handler, schema *ftonic.Adapter, require middlewares.Require) {
	group := router.Group("/roles")

	ftonic.For[apidoc.Pagination, response.Page[RoleResponse]](schema).
		GET(group, "/", require("roles", "read"), h.List, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список ролей",
			Description: "Требует roles:read.",
			Tags:        []string{"Роли"},
		}))

	ftonic.For[ListByServiceRequest, response.Page[RoleResponse]](schema).
		GET(group, "/service/:service_id", require("roles", "read"), h.ListByServiceID, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список ролей по service_id",
			Description: "Требует roles:read.",
			Tags:        []string{"Роли"},
		}))

	ftonic.For[AddPermissionRequest, RoleResponse](schema).
		POST(group, "/perm/add", require("roles.perm", "edit"), h.AddPermission, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Добавить разрешение для роли",
			Description: "Требует roles.perm:edit.",
			Tags:        []string{"Роли"},
		}))

	ftonic.For[AddPermissionRequest, RoleResponse](schema).
		POST(group, "/perm/remove", require("roles.perm", "edit"), h.RemovePermission, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Удалить разрешение у роли",
			Description: "Требует roles.perm:edit.",
			Tags:        []string{"Роли"},
		}))

	ftonic.For[UpsertRequest, RoleResponse](schema).
		POST(group, "/create", require("roles", "create"), h.Create, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Создать роль",
			Description: "Требует roles:create.",
			Tags:        []string{"Роли"},
		}))

	ftonic.For[UpdateRoleRequest, RoleResponse](schema).
		PUT(group, "/:role_id", require("roles", "update"), h.Update, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Редактировать роль",
			Description: "Требует roles:update.",
			Tags:        []string{"Роли"},
		}))

	ftonic.For[RoleIDRequest, DetailedResponse](schema).
		GET(group, "/:role_id/detailed", require("roles", "read"), h.Detailed, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Подробная информация о роли",
			Description: "Разрешения, сгруппированные по сервисам, + список пользователей с этой ролью. Требует roles:read.",
			Tags:        []string{"Роли"},
		}))

	ftonic.For[RoleIDRequest, DeleteResponse](schema).
		DELETE(group, "/:role_id", require("roles", "delete"), h.Delete, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Удалить роль",
			Description: "Требует roles:delete.",
			Tags:        []string{"Роли"},
		}))
}
