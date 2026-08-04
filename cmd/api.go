package main

import (
	"database/sql"
	"permisson/internal/auth"
	repo "permisson/internal/database/sqlc"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

type application struct {
	config config
	db     *sql.DB
	logger *log.Logger
}

type config struct {
	addr         string
	db           dbConfig
	prefix       string
	appEnv       string
	cookieDomain string
	sessionTTL   time.Duration
}

type dbConfig struct {
	dsn string
}

func (app *application) mount() *fiber.App {
	fiberApp := fiber.New(fiber.Config{})

	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://192.168.88.147:5173", "http://localhost:5173", "http://localhost:8080", "http://192.168.88.147:5176", "http://192.168.88.147:8080"},
		AllowCredentials: true,
	}))

	fiberApp.Use(logger.New(logger.Config{
		Format: "${time} | [${ip}]:${port} | ${latency} | ${status} - ${method} ${path} \n",
	}))

	v1 := fiberApp.Group("api/mag/v1")

	v1.Get("/test", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "success"})
	})

	authService := auth.NewService(repo.New(app.db), time.Hour*4)
	authHandler := auth.NewHandler(authService, app.config.appEnv == "production", app.config.cookieDomain)
	auth.RegisterRoutes(v1, authHandler)

	return fiberApp
}

func (app *application) run(f *fiber.App) error {
	return f.Listen(app.config.addr, fiber.ListenConfig{EnablePrefork: true})
}
