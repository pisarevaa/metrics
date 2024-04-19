package agent

import (
	"go.uber.org/zap"
)

func GetLogger() *zap.SugaredLogger {
	logger := zap.NewExample().Sugar()
	defer logger.Sync() //nolint:errcheck // ignore check
	return logger
}
