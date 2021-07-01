package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(input string) (string, error) {
	var (
		hash = sha256.New()
	)

	if _, err := hash.Write([]byte(input)); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
