package server

import (
	"net/http"

	"go.uber.org/zap"
)

func GetLogger() *zap.SugaredLogger {
	logger := zap.NewExample().Sugar()
	defer logger.Sync() //nolint:errcheck // ignore check
	return logger
}

type responseData struct {
	status int
	size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
