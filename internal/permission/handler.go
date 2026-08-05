package permission

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

	permissions, err := h.service.List(c.Context(), page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}

	return c.JSON(permissions)
}

// ListByServiceID — GET /permissions/:service_id. Аналог get_permissions_by_service_id.
func (h *Handler) ListByServiceID(c fiber.Ctx) error {
	serviceID := c.Params("service_id")
	page := query.QueryInt(c, "page", 1)
	limit := query.QueryInt(c, "limit", 10)

	items, err := h.service.ListByServiceID(c.Context(), serviceID, page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(items)
}

// Create — POST /permissions/create. Аналог create_permission.
func (h *Handler) Create(c fiber.Ctx) error {
	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	perm, err := h.service.Create(c.Context(), req)
	if err != nil {
		if errors.Is(err, ErrCodeExists) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.Status(fiber.StatusCreated).JSON(perm)
}

// Update — PUT /permissions/:permission_id. Аналог update_permission.
func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("permission_id")

	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	perm, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, ErrCodeExists):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
		}
	}
	return c.JSON(perm)
}

// Delete — DELETE /permissions/:permission_id. Аналог delete_permission.
func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("permission_id")

	if err := h.service.Delete(c.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(DeleteResponse{Message: "Разрешение успешно удалено"})
}
