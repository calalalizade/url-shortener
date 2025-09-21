package shortener

import (
	"crypto/sha256"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 7

func GenerateCodeFromURL(url string) string {
	hash := sha256.Sum256([]byte(url))

	code := make([]byte, codeLength)
	for i := range codeLength {
		code[i] = alphabet[int(hash[i])%len(alphabet)]
	}

	return string(code)
}
