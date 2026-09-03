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

// Logout удаляет сессию пользователя и очищает cookie "session".
func (h *Handler) Logout(c fiber.Ctx) error {
	sessionToken := c.Cookies("session")

	err := h.service.Logout(c.Context(), sessionToken)
	if err != nil {
		if errors.Is(err, ErrSessionTokenMissing) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		Value:    "",
		MaxAge:   -1,
		Secure:   h.cookieSecure,
		Domain:   h.cookieDomain,
	})

	return c.JSON(LogoutResponse{Message: "Вы успешно вышли из системы"})
}

const forgotPasswordMessage = "Если такой пользователь существует, письмо с инструкциями отправлено"

func (h *Handler) ForgotPassword(c fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Ошибку сюда бы не пропустили только сбои самой БД — Service.ForgotPassword
	// уже гасит "пользователь не найден"/SMTP-сбои внутри себя, ответ клиенту
	// в любом случае один и тот же (см. комментарий в auth.Service.ForgotPassword).
	if err := h.service.ForgotPassword(c.Context(), req.Username); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}

	return c.JSON(ForgotPasswordResponse{Message: forgotPasswordMessage})
}

func (h *Handler) ResetPassword(c fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невалидное тело запроса")
	}
	if err := req.Validate(); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := h.service.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		if errors.Is(err, ErrInvalidResetToken) {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Внутренняя ошибка")
	}

	return c.JSON(ResetPasswordResponse{Message: "Пароль изменён"})
}
