package invite

import (
	middlewares "permisson/internal/middleware"
	"permisson/internal/user"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core/docs"
	"github.com/gofiber/fiber/v3"
)

// Статический путь /invites/code/:code регистрируется раньше /invites/:invite_id
// — иначе Fiber перехватит его параметром (см. тот же приём в user/router.go).
func RegisterRoutes(router fiber.Router, h *Handler, schema *ftonic.Adapter, require middlewares.Require) {
	group := router.Group("/invites")

	ftonic.For[CodeRequest, ValidateResponse](schema).
		GET(group, "/code/:code", h.ValidateCode, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Проверить код приглашения",
			Description: "Публичный эндпоинт, без побочных эффектов. Используется фронтом перед показом формы регистрации.",
			Tags:        []string{"Приглашения"},
		}))

	ftonic.For[AcceptRequest, user.UserResponse](schema).
		POST(group, "/code/:code/accept", h.Accept, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Принять приглашение и зарегистрироваться",
			Description: "Публичный эндпоинт. Создаёт пользователя и сразу авторизует (ставит cookie session).",
			Tags:        []string{"Приглашения"},
		}))

	ftonic.For[CreateRequest, InviteResponse](schema).
		POST(group, "/", require("invites", "create"), h.Create, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Создать приглашение",
			Description: "Требует invites:create. При вызове по сервисному API-ключу created_by остаётся пустым.",
			Tags:        []string{"Приглашения"},
		}))

	ftonic.For[ListRequest, []InviteResponse](schema).
		GET(group, "/", require("invites", "read"), h.List, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить список приглашений",
			Description: "Без пагинации, требует invites:read.",
			Tags:        []string{"Приглашения"},
		}))

	ftonic.For[InviteIDRequest, InviteResponse](schema).
		GET(group, "/:invite_id", require("invites", "read"), h.GetByID, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Получить приглашение по id",
			Description: "Требует invites:read.",
			Tags:        []string{"Приглашения"},
		}))

	ftonic.For[InviteIDRequest, InviteResponse](schema).
		POST(group, "/:invite_id/revoke", require("invites", "revoke"), h.Revoke, ftonic.WithOperation(docs.OperationObject{
			Summary:     "Отозвать приглашение",
			Description: "Требует invites:revoke. Нельзя отозвать уже использованное или уже отозванное приглашение.",
			Tags:        []string{"Приглашения"},
		}))
}
