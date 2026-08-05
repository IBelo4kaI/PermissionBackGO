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
	return &Handler{service: service}
}

// List — GET /roles. Аналог role_routes.get_roles.
func (h *Handler) List(c fiber.Ctx) error {
	page := query.QueryInt(c, "page", 1)
	limit := query.QueryInt(c, "limit", 10)

	result, err := h.service.List(c.Context(), page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.JSON(result)
}

// ListByServiceID — GET /roles/service/:service_id. Аналог get_roles_by_service_id.
func (h *Handler) ListByServiceID(c fiber.Ctx) error {
	serviceID := c.Params("service_id")
	page := query.QueryInt(c, "page", 1)
	limit := query.QueryInt(c, "limit", 10)

	result, err := h.service.ListByServiceID(c.Context(), serviceID, page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.JSON(result)
}

// Create — POST /roles/create. Аналог create_role.
func (h *Handler) Create(c fiber.Ctx) error {
	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.service.Create(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

// Update — PUT /roles/:role_id. Аналог update_role.
func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("role_id")

	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.JSON(result)
}

// Delete — DELETE /roles/:role_id. Аналог delete_role.
func (h *Handler) Delete(c fiber.Ctx) error {
	id := c.Params("role_id")

	if err := h.service.Delete(c.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.JSON(DeleteResponse{Message: "Роль успешно удалена"})
}

// AddPermission — POST /roles/perm/add. Аналог role_add.
func (h *Handler) AddPermission(c fiber.Ctx) error {
	var req AddPermissionRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "невалидное тело запроса")
	}

	result, err := h.service.AddPermission(c.Context(), req.RoleID, req.PermID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound), errors.Is(err, ErrPermissionNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, ErrPermissionAlreadyAssigned):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
		}
	}
	return c.JSON(result)
}

// RemovePermission — POST /roles/perm/remove. Аналог role_remove.
func (h *Handler) RemovePermission(c fiber.Ctx) error {
	var req AddPermissionRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "невалидное тело запроса")
	}

	result, err := h.service.RemovePermission(c.Context(), req.RoleID, req.PermID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound), errors.Is(err, ErrPermissionNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		case errors.Is(err, ErrPermissionNotAssigned):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
		}
	}
	return c.JSON(result)
}

// Detailed — GET /roles/:role_id/detailed. Аналог get_role_detailed.
func (h *Handler) Detailed(c fiber.Ctx) error {
	id := c.Params("role_id")

	result, err := h.service.Detailed(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.JSON(result)
}
