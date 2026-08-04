package passwordhash

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltLen    = 32
	keyLen     = 32 // совпадает с размером digest sha256 — как дефолтный dklen в hashlib.pbkdf2_hmac
	iterations = 100_000
)

// Hash повторяет формат Python-версии: hex(salt) + hex(key), итого 128 hex-символов.
func Hash(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	key := pbkdf2.Key([]byte(password), salt, iterations, keyLen, sha256.New)

	return hex.EncodeToString(salt) + hex.EncodeToString(key), nil
}

// Verify сверяет пароль с хэшем, сохранённым в формате Hash.
func Verify(storedPassword, providedPassword string) bool {
	if len(storedPassword) != (saltLen+keyLen)*2 {
		return false
	}

	saltHex := storedPassword[:saltLen*2]
	storedKeyHex := storedPassword[saltLen*2:]

	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false
	}

	providedKey := pbkdf2.Key([]byte(providedPassword), salt, iterations, keyLen, sha256.New)
	providedKeyHex := hex.EncodeToString(providedKey)

	// В отличие от Python-версии (обычное ==), тут сравнение константного времени —
	// защита от timing-атак, поведение по сути идентично.
	return subtle.ConstantTimeCompare([]byte(storedKeyHex), []byte(providedKeyHex)) == 1
}
