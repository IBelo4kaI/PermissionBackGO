package settings

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetSMTP(c fiber.Ctx) error {
	resp, err := h.service.Get(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}
	return c.JSON(resp)
}

func (h *Handler) UpsertSMTP(c fiber.Ctx) error {
	var req UpsertSMTPSettingsRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}

	resp, err := h.service.Upsert(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrHostRequired), errors.Is(err, ErrPortRequired),
			errors.Is(err, ErrUsernameRequired), errors.Is(err, ErrFromAddressRequired),
			errors.Is(err, ErrPasswordRequired):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		default:
			return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
		}
	}
	return c.JSON(resp)
}
