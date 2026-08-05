package main

import (
	"database/sql"
	"permisson/internal/auth"
	repo "permisson/internal/database/sqlc"
	"permisson/internal/permission"
	"permisson/internal/role"
	"permisson/internal/service"
	"permisson/internal/user"
	"time"

	ftonic "github.com/TickLabVN/tonic/adapters/fiber"
	"github.com/TickLabVN/tonic/core"
	"github.com/TickLabVN/tonic/core/docs"
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
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8080"},
		AllowCredentials: true,
	}))

	fiberApp.Use(logger.New(logger.Config{
		Format: "${time} | [${ip}]:${port} | ${latency} | ${status} - ${method} ${path} \n",
	}))

	schema := ftonic.New(&docs.OpenApi{
		OpenAPI: docs.VERSION,
		Info: docs.InfoObject{
			Title: "Сервис аутентификации и контроля разрешений",
		},
	})

	v1 := fiberApp.Group("/api/as")

	authService := auth.NewService(repo.New(app.db), time.Hour*4)
	authHandler := auth.NewHandler(authService, app.config.appEnv == "production", app.config.cookieDomain)
	auth.RegisterRoutes(v1, authHandler, schema)

	permissionService := permission.NewService(repo.New(app.db))
	permissionHandler := permission.NewHandler(permissionService)
	permission.RegisterRoutes(v1, permissionHandler, schema)

	serviceService := service.NewService(repo.New(app.db))
	serviceHandler := service.NewHandler(serviceService, authService)
	service.RegisterRoutes(v1, serviceHandler, schema)

	userService := user.NewService(repo.New(app.db))
	userHandler := user.NewHandler(userService, authService)
	user.RegisterRoutes(v1, userHandler, schema)

	roleService := role.NewService(repo.New(app.db))
	roleHandler := role.NewHandler(roleService)
	role.RegisterRoutes(v1, roleHandler, schema)

	schema.UIHandle(fiberApp, "/docs", core.SwaggerUI)

	return fiberApp
}

func (app *application) run(f *fiber.App) error {
	return f.Listen(app.config.addr, fiber.ListenConfig{EnablePrefork: true})
}
