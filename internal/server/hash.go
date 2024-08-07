package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
)

// Получение хеша тела запроса с добавлением соли key.
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

// Мидлвар про проверке хеша запроса.
func (s *Handler) HashCheckMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hashInHeader := r.Header.Get("Hash")

		var buf bytes.Buffer
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		payload := buf.Bytes()
		s.Logger.Error("payload:", string(payload))
		hash, errHash := GetBodyHash(payload, s.Config.Key)
		if errHash != nil {
			s.Logger.Error(errHash)
			http.Error(w, errHash.Error(), http.StatusBadRequest)
			return
		}
		if hashInHeader != "none" && hashInHeader != "" {
			if hashInHeader != hash {
				s.Logger.Error("Hash mismatch:", hashInHeader, "!=", hash)
				http.Error(w, "Hash mismatch", http.StatusBadRequest)
				return
			}
		}

		r.Body = io.NopCloser(bytes.NewReader(payload))
		h.ServeHTTP(w, r)
	})
}
