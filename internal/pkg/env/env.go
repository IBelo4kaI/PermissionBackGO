package env

import (
	"encoding/base64"
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

// GetCorporateDbString — DSN общей (corporate) БД: профили/оргструктура,
// которыми этот сервис не управляет (миграции туда не применяются, см.
// internal/database/corporate_schema), но пишет напрямую при принятии инвайта.
func (e *Env) GetCorporateDbString() string {
	return os.Getenv("CORPORATE_DB_STRING")
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

// GetAppBaseURL — базовый URL фронта, нужен только чтобы собрать ссылку на
// сброс пароля в письме (APP_BASE_URL + "/reset-password/" + rawToken).
func (e *Env) GetAppBaseURL() string {
	return os.Getenv("APP_BASE_URL")
}

// GetSettingsEncryptionKey — ключ AES-256 (32 байта, base64) для
// internal/pkg/secretcrypt. Используется только для SMTP-пароля в
// internal/settings — единственного секрета в проекте, который нужно уметь
// расшифровать обратно, а не только хэшировать/сравнивать.
func (e *Env) GetSettingsEncryptionKey() ([]byte, error) {
	raw := os.Getenv("SETTINGS_ENCRYPTION_KEY")
	if raw == "" {
		return nil, fmt.Errorf("SETTINGS_ENCRYPTION_KEY обязателен")
	}

	key, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("SETTINGS_ENCRYPTION_KEY должен быть в base64: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("SETTINGS_ENCRYPTION_KEY должен декодироваться в 32 байта (AES-256), получено %d", len(key))
	}

	return key, nil
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
