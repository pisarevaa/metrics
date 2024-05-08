package server

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func createMerticsTable(dbpool *pgxpool.Pool) error {
	_, err := dbpool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS metrics (
			"id" VARCHAR(250) PRIMARY KEY,
			"type" VARCHAR(50) NOT NULL,
			"delta" INTEGER NULL,
			"value" DOUBLE PRECISION NULL,
			"updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
		) 
	`)
	if err != nil {
		return err
	}
	return nil
}

func restoreMetricsFromDB(dbpool *pgxpool.Pool, storage *MemStorage) error {
	rows, err := dbpool.Query(context.Background(), "SELECT id, type, delta, value FROM metrics")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var m Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			return err
		}
		storage.Store(m)
	}
	return nil
}

func ConnectDB(config Config, logger *zap.SugaredLogger) *pgxpool.Pool {
	if config.DatabaseDSN == "" {
		return nil
	}
	dbpool, err := pgxpool.New(context.Background(), config.DatabaseDSN)
	if err != nil {
		logger.Error("Unable to create connection pool: %v", err)
		return nil
	}
	err = createMerticsTable(dbpool)
	if err != nil {
		logger.Error("Unable to create table metrics: %v", err)
		return nil
	}
	return dbpool
}
