package invite

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"permisson/internal/auth"
	"permisson/internal/pkg/query"
	"permisson/internal/user"
)

// Handler держит *auth.Service по двум причинам: определить created_by (см.
// Create) и автологин сразу после Accept (как auth.Handler.Login).
type Handler struct {
	service      *Service
	auth         *auth.Service
	cookieSecure bool
	cookieDomain string
}

func NewHandler(service *Service, authService *auth.Service, cookieSecure bool, cookieDomain string) *Handler {
	return &Handler{service: service, auth: authService, cookieSecure: cookieSecure, cookieDomain: cookieDomain}
}

// Create: created_by = callerInfo.UserID — для пользовательской сессии это
// id из сессии, для сервисного API-ключа CallerFrom его не заполняет, так
// что поле остаётся "" и Service.Create кладёт NULL (см.
// schema/004_invites_created_by_nullable.sql) — специального разбора по
// типу caller'а тут не нужно.
func (h *Handler) Create(c fiber.Ctx) error {
	callerInfo, err := h.auth.CallerFrom(c.Context(), c.Cookies("session"))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	var req CreateRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}

	invite, err := h.service.Create(c.Context(), req, callerInfo.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.Status(fiber.StatusCreated).JSON(invite)
}

func (h *Handler) List(c fiber.Ctx) error {
	search := query.QuerySearch(c)
	companyID := query.QueryStringPtr(c, "company_id")
	departmentID := query.QueryStringPtr(c, "department_id")
	positionID := query.QueryStringPtr(c, "position_id")
	sortBy, sortDir := query.QuerySort(c, SortableColumns, DefaultSortColumn)

	items, err := h.service.List(c.Context(), search, companyID, departmentID, positionID, sortBy, sortDir)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(items)
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id := c.Params("invite_id")

	invite, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(invite)
}

func (h *Handler) Revoke(c fiber.Ctx) error {
	id := c.Params("invite_id")

	invite, err := h.service.Revoke(c.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, ErrInviteUsed), errors.Is(err, ErrInviteRevoked):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
		}
	}
	return c.JSON(invite)
}

// ValidateCode — публичный, без require: фронт проверяет ссылку до показа
// формы регистрации.
func (h *Handler) ValidateCode(c fiber.Ctx) error {
	code := c.Params("code")

	resp, err := h.service.ValidateCode(c.Context(), code)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(resp)
}

// Accept — публичный, без require: сам факт валидного кода — это и есть
// авторизация на регистрацию. После успешного создания аккаунта сразу
// логинит (как auth.Handler.Login) и ставит cookie "session".
func (h *Handler) Accept(c fiber.Ctx) error {
	var req AcceptRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}
	req.Code = c.Params("code")

	newUser, err := h.service.Accept(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrCodeNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, ErrInviteUsed), errors.Is(err, ErrInviteRevoked), errors.Is(err, ErrInviteExpired),
			errors.Is(err, user.ErrGenderNotFound), errors.Is(err, user.ErrUsernameExists),
			errors.Is(err, user.ErrNameRequired), errors.Is(err, user.ErrSurnameRequired),
			errors.Is(err, user.ErrUsernameRequired), errors.Is(err, user.ErrPasswordRequired),
			errors.Is(err, user.ErrBirthdayRequired), errors.Is(err, user.ErrGenderRequired):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
		}
	}

	loginResult, err := h.auth.Login(c.Context(), auth.LoginRequest{Login: req.Username, Password: req.Password})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		Value:    loginResult.Token,
		Expires:  loginResult.ExpiresAt,
		Secure:   h.cookieSecure,
		Domain:   h.cookieDomain,
	})

	return c.Status(fiber.StatusCreated).JSON(newUser)
}
