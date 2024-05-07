package server

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func ConnectDB(config Config, logger *zap.SugaredLogger) *pgxpool.Pool {
	if config.DatabaseDSN == "" {
		return nil
	}
	fmt.Println("config.DatabaseDSN", config.DatabaseDSN)
	dbpool, err := pgxpool.New(context.Background(), config.DatabaseDSN)
	if err != nil {
		logger.Error("Unable to create connection pool: %v", err)
		return nil
	}
	return dbpool
}
