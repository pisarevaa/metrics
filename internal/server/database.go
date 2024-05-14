package server

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type MetricsModel interface {
	IsExist() (check bool)
	Ping(ctx context.Context) (err error)
	RestoreMetricsFromDB(storage *MemStorage) (err error)
	InsertRowsIntoDB(ctx context.Context, metrics []Metrics) (err error)
	InsertRowIntoDB(ctx context.Context, metric Metrics, now time.Time) (err error)
}

type DBPool struct {
	*pgxpool.Pool
}

func (dbpool *DBPool) IsExist() bool {
	return dbpool != nil
}

func (dbpool *DBPool) Ping(ctx context.Context) error {
	var one int
	err := dbpool.QueryRow(ctx, "select 1").Scan(&one)
	if err != nil {
		return err
	}
	return nil
}

func (dbpool *DBPool) RestoreMetricsFromDB(storage *MemStorage) error {
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

func (dbpool *DBPool) InsertRowsIntoDB(ctx context.Context, metrics []Metrics) error {
	tx, errTx := dbpool.Begin(ctx)
	if errTx != nil {
		return errTx
	}
	now := time.Now()
	for _, metric := range metrics {
		err := dbpool.InsertRowIntoDB(
			ctx,
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

func (dbpool *DBPool) InsertRowIntoDB(ctx context.Context, metric Metrics, now time.Time) error {
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

func NewDBPool(config Config, logger *zap.SugaredLogger) *DBPool {
	if config.DatabaseDSN == "" {
		return nil
	}
	dbpool, err := pgxpool.New(context.Background(), config.DatabaseDSN)
	if err != nil {
		logger.Error("Unable to create connection pool: %v", err)
		return nil
	}
	m, err := migrate.New(
		"file://migrations",
		config.DatabaseDSN)
	if err != nil {
		logger.Error("Unable to migrate tables: ", err)
		return nil
	}
	err = m.Up()
	if err != nil {
		logger.Error("Unable to migrate tables: ", err)
		return nil
	}
	db := &DBPool{dbpool}
	return db
}
