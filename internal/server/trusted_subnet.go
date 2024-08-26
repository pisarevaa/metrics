package server

import (
	"net/http"
	"strings"
)

// Мидлвар про проверке IP запроса.
func (s *Handler) IPCheckMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			http.Error(w, "X-Real-IP header is empty", http.StatusBadRequest)
			return
		}
		if !strings.Contains(s.Config.TrustedSubnet, ip) {
			http.Error(w, "IP is not trusted", http.StatusBadRequest)
			return
		}
		h.ServeHTTP(w, r)
	})
}
