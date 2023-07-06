package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

func VerifyAPIKey(keyPlainText string) bool {
	hash := sha256.Sum256([]byte(keyPlainText))
	if hex.EncodeToString(hash[:]) == os.Getenv("API_HASH") {
		return true
	}
	return false
}
