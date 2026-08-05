package service

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
	return &Handler{service: service, auth: authService}
}

// List — GET /services. Аналог service_routers.get_all.
func (h *Handler) List(c fiber.Ctx) error {
	page := query.QueryInt(c, "page", 1)
	limit := query.QueryInt(c, "limit", 10)

	result, err := h.service.List(c.Context(), page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.JSON(result)
}

// Create — POST /services/create. Аналог service_routers.create.
func (h *Handler) Create(c fiber.Ctx) error {
	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	svc, err := h.service.Create(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.Status(fiber.StatusCreated).JSON(svc)
}

// Update — PUT /services/:service_id. Аналог service_routers.update.
func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("service_id")

	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	svc, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.JSON(svc)
}

// ListUserAccessible — GET /services/user-accessible. Аналог
// get_user_accessible_services: требует ТОЛЬКО пользовательскую сессию
// (Depends(get_session)), в отличие от остальных роутов сущности, которые
// в Python защищены require_permission.
func (h *Handler) ListUserAccessible(c fiber.Ctx) error {
	sessionToken := c.Cookies("session")

	userID, err := h.auth.SessionUserID(c.Context(), sessionToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "невалидная сессия")
	}

	items, err := h.service.ListAccessibleForUser(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.JSON(items)
}

// IssueAPIKey — POST /services/:service_id/api-key. Аналог issue_api_key.
func (h *Handler) IssueAPIKey(c fiber.Ctx) error {
	id := c.Params("service_id")

	key, err := h.service.IssueAPIKey(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.JSON(key)
}

// RevokeAPIKey — DELETE /services/:service_id/api-key. Аналог revoke_api_key.
func (h *Handler) RevokeAPIKey(c fiber.Ctx) error {
	id := c.Params("service_id")

	if err := h.service.RevokeAPIKey(c.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "внутренняя ошибка")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
