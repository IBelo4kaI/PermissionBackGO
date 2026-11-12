package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v3"
	"permisson/internal/database"
)

var tenantID = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// TENANTS is an explicit allowlist. Forwarded host headers are deliberately ignored.
func runTenants(base config, logger *log.Logger) error {
	router := fiber.New()
	handlers := map[string]*fiber.App{}
	names := map[string]bool{}
	databases := map[string]bool{}
	corporateDatabases := map[string]bool{}
	for _, rawID := range strings.Split(os.Getenv("TENANTS"), ",") {
		id := strings.TrimSpace(rawID)
		if !tenantID.MatchString(id) || names[id] {
			return fmt.Errorf("invalid or duplicate tenant %q", id)
		}
		names[id] = true
		get := func(k string) string { return strings.TrimSpace(os.Getenv("TENANT_" + id + "_" + k)) }
		host := strings.ToLower(get("HOST"))
		if host == "" || strings.ContainsAny(host, "/: \t") || handlers[host] != nil {
			return fmt.Errorf("tenant %s: invalid or duplicate HOST", id)
		}
		dsn, corporate := get("DB_STRING"), get("CORPORATE_DB_STRING")
		parsed, err := mysql.ParseDSN(dsn)
		if err != nil || parsed == nil || parsed.DBName == "" {
			return fmt.Errorf("tenant %s: invalid DB_STRING", id)
		}
		dbKey := parsed.Net + ":" + parsed.Addr + "/" + parsed.DBName
		if databases[dbKey] {
			return fmt.Errorf("tenant %s: authentication database is already assigned", id)
		}
		databases[dbKey] = true
		corpParsed, err := mysql.ParseDSN(corporate)
		if err != nil || corpParsed == nil || corpParsed.DBName == "" {
			return fmt.Errorf("tenant %s: invalid CORPORATE_DB_STRING", id)
		}
		corpKey := corpParsed.Net + ":" + corpParsed.Addr + "/" + corpParsed.DBName
		if corporateDatabases[corpKey] {
			return fmt.Errorf("tenant %s: corporate database is already assigned", id)
		}
		corporateDatabases[corpKey] = true
		baseURL := get("APP_BASE_URL")
		u, err := url.Parse(baseURL)
		if err != nil || u.Hostname() == "" || (u.Scheme != "https" && u.Scheme != "http") {
			return fmt.Errorf("tenant %s: invalid APP_BASE_URL", id)
		}
		domain := strings.TrimPrefix(strings.ToLower(get("COOKIE_DOMAIN")), ".")
		if domain != "" && host != domain && !strings.HasSuffix(host, "."+domain) {
			return fmt.Errorf("tenant %s: COOKIE_DOMAIN does not match HOST", id)
		}
		origins := []string{}
		for _, origin := range strings.Split(get("ALLOWED_ORIGINS"), ",") {
			if origin = strings.TrimSpace(origin); origin != "" {
				origins = append(origins, origin)
			}
		}
		if len(origins) == 0 {
			origins = []string{u.Scheme + "://" + u.Host}
		}
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return fmt.Errorf("tenant %s: %w", id, err)
		}
		defer db.Close()
		corp, err := sql.Open("mysql", corporate)
		if err != nil {
			return fmt.Errorf("tenant %s corporate DB: %w", id, err)
		}
		defer corp.Close()
		if err := database.RunMigrations(context.Background(), db); err != nil {
			return fmt.Errorf("tenant %s migrations: %w", id, err)
		}
		cfg := base
		cfg.cookieName = "session_" + strings.ToLower(id)
		cfg.cookieDomain = domain
		cfg.appBaseURL = baseURL
		cfg.allowedOrigins = origins
		app := application{config: cfg, db: db, corporateDB: corp, logger: logger}
		handlers[host] = app.mount()
	}
	router.Use(func(c fiber.Ctx) error {
		host := strings.ToLower(c.Get("Host"))
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		app := handlers[host]
		if app == nil {
			return c.SendStatus(fiber.StatusMisdirectedRequest)
		}
		app.Handler()(c.RequestCtx())
		return nil
	})
	logger.Info("multi-tenant mode: gRPC is disabled; HTTP tenants are selected by Host")
	return router.Listen(base.addr)
}
