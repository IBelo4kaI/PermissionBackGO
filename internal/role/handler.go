package role

import (
	"errors"
	"permisson/internal/pkg/query"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) List(c fiber.Ctx) error {
	page := query.QueryInt(c, "page", 1)
	limit := query.QueryInt(c, "limit", 10)

	roles, err := h.service.List(c.Context(), page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(roles)
}

func (h *Handler) ListByServiceID(c fiber.Ctx) error {
	serviceID := c.Params("service_id")
	page := query.QueryInt(c, "page", 1)
	limit := query.QueryInt(c, "limit", 10)

	roles, err := h.service.ListByServiceID(c.Context(), serviceID, page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(roles)
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	role, err := h.service.Create(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.Status(fiber.StatusCreated).JSON(role)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("role_id")

	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	role, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(role)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("role_id")

	if err := h.service.Delete(c.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(DeleteResponse{Message: "Роль успешно удалена"})
}

func (h *Handler) AddPermission(c fiber.Ctx) error {
	var req AddPermissionRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}

	role, err := h.service.AddPermission(c.Context(), req.RoleID, req.PermID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound), errors.Is(err, ErrPermissionNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, ErrPermissionAlreadyAssigned):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
		}
	}
	return c.JSON(role)
}

func (h *Handler) RemovePermission(c fiber.Ctx) error {
	var req AddPermissionRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}

	role, err := h.service.RemovePermission(c.Context(), req.RoleID, req.PermID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound), errors.Is(err, ErrPermissionNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, ErrPermissionNotAssigned):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
		}
	}
	return c.JSON(role)
}

func (h *Handler) Detailed(c fiber.Ctx) error {
	id := c.Params("role_id")

	detailed, err := h.service.Detailed(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(detailed)
}
