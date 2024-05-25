package utils

import (
	"crypto/sha256"
	"encoding/base64"
)

func GetBodyHash(payload []byte, key string) (string, error) {
	payload = append(payload, []byte(key)...)
	h := sha256.New()
	_, err := h.Write(payload)
	if err != nil {
		return "", err
	}
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha, nil
}
