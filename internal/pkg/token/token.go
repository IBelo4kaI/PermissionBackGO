package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// Generate возвращает случайный токен: 32 случайных байта в hex (64 символа) —
// аналог Python secrets.token_hex(32).
func Generate() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Hash возвращает sha256-хэш токена в hex — аналог Python
// hashlib.sha256(token.encode()).hexdigest(). В БД хранится именно хэш, не сам токен.
func Hash(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
