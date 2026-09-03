// Package secretcrypt — симметричное шифрование секретов, которые нужно
// уметь прочитать обратно (не хэш): единственный такой случай в проекте —
// SMTP-пароль в internal/settings, который реально нужен в открытом виде,
// чтобы аутентифицироваться на почтовом сервере.
package secretcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrInvalidKey         = errors.New("ключ шифрования должен быть 32 байта (AES-256)")
	ErrCiphertextTooShort = errors.New("некорректный шифротекст")
)

func newGCM(key []byte) (cipher.AEAD, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKey
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

// Encrypt шифрует plaintext AES-256-GCM. Результат — base64(nonce || ciphertext),
// готов для хранения в текстовом столбце БД.
func Encrypt(key []byte, plaintext string) (string, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt — обратная операция Encrypt.
func Decrypt(key []byte, encoded string) (string, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrCiphertextTooShort
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
