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
	keyLen     = 32
	iterations = 100_000
)

func Hash(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	key := pbkdf2.Key([]byte(password), salt, iterations, keyLen, sha256.New)

	return hex.EncodeToString(salt) + hex.EncodeToString(key), nil
}

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

	return subtle.ConstantTimeCompare([]byte(storedKeyHex), []byte(providedKeyHex)) == 1
}
