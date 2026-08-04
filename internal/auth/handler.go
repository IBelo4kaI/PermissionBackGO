package auth

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	service      *Service
	cookieSecure bool
	cookieDomain string
}

func NewHandler(service *Service, cookieSecure bool, cookieDomain string) *Handler {
	return &Handler{
		service:      service,
		cookieSecure: cookieSecure,
		cookieDomain: cookieDomain,
	}
}

func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}

	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	result, err := h.service.Login(c.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return fiber.NewError(fiber.StatusNotFound, ErrInvalidCredentials.Error())
		}
		if errors.Is(err, ErrUserInactive) {
			return fiber.NewError(fiber.StatusNotFound, ErrUserInactive.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		Value:    result.Token,
		Expires:  result.ExpiresAt,
		Secure:   h.cookieSecure,
		Domain:   h.cookieDomain,
	})

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) ValidateSession(c fiber.Ctx) error {
	sessionToken := c.Cookies("session")

	valid, err := h.service.ValidateSession(c.Context(), sessionToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}

	return c.JSON(ValidateSessionResponse{Valid: valid})
}
