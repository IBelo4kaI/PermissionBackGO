package auth

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, h *Handler) {
	group := router.Group("/auth")
	group.Post("/login", h.Login)
	group.Post("/validate-session", h.ValidateSession)
}
