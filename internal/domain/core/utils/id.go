package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateUniqueID создает случайный 16-байтовый идентификатор (32-символьная строка)
func GenerateUniqueID() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
