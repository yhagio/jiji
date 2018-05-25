package utils

import (
	"crypto/rand"
	"encoding/base64"
)

const TokenBytes = 32

// Generate token of a predetermined byte size.
func GenerateToken() (string, error) {
	return generateEncodedString(TokenBytes)
}

// Generate n random bytes.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Generate a byte slice of size nBytes and then
// return a string that is the base64 URL encoded version
// of that byte slice
func generateEncodedString(nBytes int) (string, error) {
	b, err := generateRandomBytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
