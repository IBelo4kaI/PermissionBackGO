package auth

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, h *Handler) {
	// group := router.Group("/auth")

	// group.Post("/register", h.Register)
	// group.Post("/login", h.Login)
	// group.Post("/refresh", h.Refresh)
	// group.Post("/logout", h.Logout)
	// group.Get("/me", Middleware(jwtSecret), h.Me)
	// group.Patch("/password", Middleware(jwtSecret), h.ChangePassword)
}
