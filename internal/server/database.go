package server

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func CreateMerticsTable(dbpool *pgxpool.Pool) error {
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

func RestoreMetricsFromDB(dbpool *pgxpool.Pool, storage *MemStorage) error {
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

func InsertRowsIntoDDB(ctx context.Context, dbpool *pgxpool.Pool, metrics []Metrics) error {
	tx, errTx := dbpool.Begin(ctx)
	if errTx != nil {
		return errTx
	}
	now := time.Now()
	for _, metric := range metrics {
		err := InsertRowIntoDDB(
			ctx,
			dbpool,
			metric,
			now,
		)
		if err != nil {
			return err
		}
	}
	errTx = tx.Commit(ctx)
	if errTx != nil {
		return errTx
	}
	return nil
}

func InsertRowsIntoDDWithRetry(ctx context.Context, dbpool *pgxpool.Pool, metrics []Metrics) error {
	retries := 3
	timeouts := map[int]int{1: 5, 2: 3, 3: 1} //nolint:gomnd // omit
	for retries > 0 {
		err := InsertRowsIntoDDB(
			ctx,
			dbpool,
			metrics,
		)
		if err != nil { //nolint:nestif // omit
			if strings.Contains(err.Error(), "failed to connect") {
				time.Sleep(time.Duration(timeouts[retries]) * time.Second)
				retries--
				if retries == 0 {
					return err
				}
			} else {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func InsertRowIntoDDB(ctx context.Context, dbpool *pgxpool.Pool, metric Metrics, now time.Time) error {
	_, err := dbpool.Exec(ctx, `
			INSERT INTO metrics (id, type, delta, value, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO UPDATE
			SET
			type = excluded.type,
			delta = excluded.delta,
			value = excluded.value,
			updated_at = excluded.updated_at
		`, metric.ID, metric.MType, metric.Delta, metric.Value, now)
	if err != nil {
		return err
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
	err = CreateMerticsTable(dbpool)
	if err != nil {
		logger.Error("Unable to create table metrics: %v", err)
		return nil
	}
	return dbpool
}
