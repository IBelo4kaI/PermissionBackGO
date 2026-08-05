package user

import (
	"permisson/internal/permission"
	"permisson/internal/pkg/apidoc"
	"permisson/internal/pkg/response"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

// require_permission пока не подключён — роуты работают без проверки прав
// до появления middleware.RequirePermission.
//
// Статические пути (/service/:service_id, /all, /me, /me/permissions/:service_id)
// регистрируются раньше /:user_id, иначе Fiber перехватит их параметром.
func RegisterRoutes(router fiber.Router, h *Handler, schema *ftonic.Adapter) {
	group := router.Group("/users")

	ftonic.For[apidoc.Pagination, response.Page[UserResponse]](schema).
		GET(group, "/", h.List, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список пользователей",
			Description: "Требует users:read_all.",
			Tags:        []string{"Пользователи"},
		}))

	ftonic.For[ListByServiceRequest, response.Page[UserResponse]](schema).
		GET(group, "/service/:service_id", h.ListByServiceID, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список пользователей по service_id",
			Description: "Требует users:read_all.",
			Tags:        []string{"Пользователи"},
		}))

	ftonic.For[apidoc.Empty, []UserResponse](schema).
		GET(group, "/all", h.ListAll, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список пользователей без пагинации",
			Description: "Требует users:read_all.",
			Tags:        []string{"Пользователи"},
		}))

	ftonic.For[apidoc.Empty, UserResponse](schema).
		GET(group, "/me", h.Me, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить текущего пользователя по сессии",
			Description: "Требует валидную пользовательскую cookie-сессию.",
			Tags:        []string{"Пользователи"},
		}))

	ftonic.For[MePermissionsRequest, []permission.Permission](schema).
		GET(group, "/me/permissions/:service_id", h.MePermissions, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список разрешений текущего пользователя для сервиса",
			Description: "Учитывает wildcard-коды вида all:all:all и <service_name>:all:all. Требует валидную пользовательскую cookie-сессию.",
			Tags:        []string{"Пользователи"},
		}))

	ftonic.For[UserIDRequest, UserResponse](schema).
		GET(group, "/:user_id", h.GetByID, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить пользователя по id",
			Description: "Требует users:read.",
			Tags:        []string{"Пользователи"},
		}))

	ftonic.For[UpdateUserRequest, UserResponse](schema).
		PUT(group, "/:user_id", h.Update, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Обновить пользователя",
			Description: "Частичное обновление: передавать можно любое подмножество полей. Требует users:update.",
			Tags:        []string{"Пользователи"},
		}))

	ftonic.For[UserIDRequest, DeleteResponse](schema).
		DELETE(group, "/:user_id", h.Delete, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Удалить пользователя",
			Description: "Требует users:delete.",
			Tags:        []string{"Пользователи"},
		}))

	ftonic.For[CreateRequest, UserResponse](schema).
		POST(group, "/create", h.Create, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Создать пользователя",
			Description: "Пароль хэшируется PBKDF2-SHA256 (как в Python-версии).",
			Tags:        []string{"Пользователи"},
		}))

	ftonic.For[RoleRequest, UserResponse](schema).
		POST(group, "/roles/add", h.AddRole, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Добавить роль пользователю",
			Description: "Требует users.roles:edit.",
			Tags:        []string{"Пользователи"},
		}))

	ftonic.For[RoleRequest, UserResponse](schema).
		POST(group, "/roles/remove", h.RemoveRole, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Удалить роль у пользователя",
			Description: "Требует users.roles:edit.",
			Tags:        []string{"Пользователи"},
		}))
}
