package user

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"permisson/internal/auth"
	"permisson/internal/pkg/query"
)

type Handler struct {
	service *Service
	auth    *auth.Service
}

func NewHandler(service *Service, authService *auth.Service) *Handler {
	return &Handler{
		service: service,
		auth:    authService,
	}
}

func (h *Handler) List(c fiber.Ctx) error {
	page := query.QueryInt(c, "page", 1)
	limit := query.QueryInt(c, "limit", 10)

	users, err := h.service.List(c.Context(), page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(users)
}

func (h *Handler) ListByServiceID(c fiber.Ctx) error {
	serviceID := c.Params("service_id")
	page := query.QueryInt(c, "page", 1)
	limit := query.QueryInt(c, "limit", 50)

	users, err := h.service.ListByServiceID(c.Context(), serviceID, page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(users)
}

func (h *Handler) ListAll(c fiber.Ctx) error {
	users, err := h.service.ListAll(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(users)
}

// Требует только пользовательскую сессию (не API-ключ), как в Python-версии.
func (h *Handler) Me(c fiber.Ctx) error {
	sessionToken := c.Cookies("session")

	userID, err := h.auth.SessionUserID(c.Context(), sessionToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Невалидная сессия")
	}

	user, err := h.service.GetByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(user)
}

func (h *Handler) MePermissions(c fiber.Ctx) error {
	sessionToken := c.Cookies("session")

	userID, err := h.auth.SessionUserID(c.Context(), sessionToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Невалидная сессия")
	}

	serviceID := c.Params("service_id")
	permissions, err := h.service.MePermissions(c.Context(), userID, serviceID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(permissions)
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id := c.Params("user_id")

	user, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(user)
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req CreateRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := h.service.Create(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrGenderNotFound), errors.Is(err, ErrUsernameExists):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
		}
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("user_id")

	var req UpdateRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}

	user, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, ErrGenderNotFound), errors.Is(err, ErrUsernameExists):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
		}
	}
	return c.JSON(user)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("user_id")

	if err := h.service.Delete(c.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(DeleteResponse{Message: "Пользователь успешно удален"})
}

func (h *Handler) AddRole(c fiber.Ctx) error {
	var req RoleRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}

	user, err := h.service.AddRole(c.Context(), req.UserID, req.RoleID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound), errors.Is(err, ErrRoleNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, ErrRoleAlreadyAssigned):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
		}
	}
	return c.JSON(user)
}

func (h *Handler) RemoveRole(c fiber.Ctx) error {
	var req RoleRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}

	user, err := h.service.RemoveRole(c.Context(), req.UserID, req.RoleID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound), errors.Is(err, ErrRoleNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, ErrRoleNotAssigned):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
		}
	}
	return c.JSON(user)
}
