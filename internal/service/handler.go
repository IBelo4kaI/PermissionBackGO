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
	return &Handler{
		service: service,
		auth:    authService,
	}
}

func (h *Handler) List(c fiber.Ctx) error {
	page := query.QueryInt(c, "page", 1)
	limit := query.QueryInt(c, "limit", 10)

	services, err := h.service.List(c.Context(), page, limit)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(services)
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	svc, err := h.service.Create(c.Context(), req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.Status(fiber.StatusCreated).JSON(svc)
}

func (h *Handler) Update(c fiber.Ctx) error {
	id := c.Params("service_id")

	var req UpsertRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	svc, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(svc)
}

// Требует только пользовательскую сессию (не API-ключ), в отличие от
// остальных роутов сущности.
func (h *Handler) ListUserAccessible(c fiber.Ctx) error {
	sessionToken := c.Cookies("session")

	userID, err := h.auth.SessionUserID(c.Context(), sessionToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Невалидная сессия")
	}

	items, err := h.service.ListAccessibleForUser(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(items)
}

func (h *Handler) IssueAPIKey(c fiber.Ctx) error {
	id := c.Params("service_id")

	key, err := h.service.IssueAPIKey(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(key)
}

func (h *Handler) RevokeAPIKey(c fiber.Ctx) error {
	id := c.Params("service_id")

	if err := h.service.RevokeAPIKey(c.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
