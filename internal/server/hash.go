package server

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
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
		hash, errHash := GetBodyHash(payload, s.Config.Key)
		if errHash != nil {
			http.Error(w, errHash.Error(), http.StatusBadRequest)
			return
		}
		if hashInHeader != hash {
			http.Error(w, "Hash mismatch", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(payload))
		h.ServeHTTP(w, r)
	})
}
