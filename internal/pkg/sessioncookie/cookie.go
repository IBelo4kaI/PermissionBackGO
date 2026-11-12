package sessioncookie

import "github.com/gofiber/fiber/v3"

const key = "tenant_session_cookie"

func Set(c fiber.Ctx, name string) { c.Locals(key, name) }
func Name(c fiber.Ctx) string {
	if name, ok := c.Locals(key).(string); ok && name != "" {
		return name
	}
	return "session"
}
