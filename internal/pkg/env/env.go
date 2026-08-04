package env

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Env struct {
}

func (e *Env) Init() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return nil
}

func (e *Env) GetDbString() string {
	return os.Getenv("DB_STRING")
}

func (e *Env) GetAddr() string {
	return os.Getenv("ADDR")
}

func (e *Env) GetAppEnv() string {
	return os.Getenv("APP_ENV")
}

func (e *Env) GetCookieDomain() string {
	return os.Getenv("COOKIE_DOMAIN")
}

func (e *Env) GetSessionTTL(fallback time.Duration) (time.Duration, error) {
	raw := os.Getenv("SESSION_TTL")
	if raw == "" {
		return fallback, nil
	}

	d, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid SESSION_TTL: %w", err)
	}

	return d, nil
}
