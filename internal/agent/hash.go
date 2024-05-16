package agent

import (
	"crypto/sha256"
)

func GetBodyHash(payload []byte, key string) (string, error) {
	payload = append(payload, []byte(key)...)
	h := sha256.New()
	_, err := h.Write(payload)
	if err != nil {
		return "", err
	}
	dst := h.Sum(nil)
	return string(dst), nil
}
