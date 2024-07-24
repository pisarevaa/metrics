package storage

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DBStorage struct {
	*pgxpool.Pool
}

// Создание БД движка хранение метрик.
func NewDBStorage(dsn string, logger *zap.SugaredLogger) *DBStorage {
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Error("Unable to create connection pool: %v", err)
		return nil
	}
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logger.Error("Unable to migrate tables: ", err)
	}
	err = m.Up()
	if err != nil {
		logger.Error("Unable to migrate tables: ", err)
	}
	db := &DBStorage{dbpool}
	return db
}

// Сохранение метрики.
func (dbpool *DBStorage) StoreMetric(ctx context.Context, metric Metrics) error {
	now := time.Now()
	_, err := dbpool.Exec(ctx, `
			INSERT INTO metrics (id, type, delta, value, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO UPDATE
			SET
			type = excluded.type,
			delta = (SELECT delta FROM metrics WHERE id = $6) + excluded.delta,
			value = excluded.value,
			updated_at = excluded.updated_at
		`, metric.ID, metric.MType, metric.Delta, metric.Value, now, metric.ID)
	if err != nil {
		return err
	}
	return nil
}

// Сохранение метрик.
func (dbpool *DBStorage) StoreMetrics(ctx context.Context, metrics []Metrics) error {
	tx, errTx := dbpool.Begin(ctx)
	if errTx != nil {
		return errTx
	}
	for _, metric := range metrics {
		err := dbpool.StoreMetric(
			ctx,
			metric,
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

// Получение метрики.
func (dbpool *DBStorage) GetMetric(ctx context.Context, id string, mtype string) (Metrics, error) {
	var m Metrics
	err := dbpool.QueryRow(ctx, "SELECT id, type, delta, value FROM metrics WHERE id = $1 and type = $2", id, mtype).
		Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
	if err != nil {
		return m, err
	}
	return m, nil
}

// Получение метрик.
func (dbpool *DBStorage) GetAllMetrics(ctx context.Context) ([]Metrics, error) {
	var metrics []Metrics
	rows, err := dbpool.Query(ctx, "SELECT id, type, delta, value FROM metrics")
	if err != nil {
		return []Metrics{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var m Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
		if err != nil {
			return []Metrics{}, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

// Пинг хранилища.
func (dbpool *DBStorage) Ping(ctx context.Context) error {
	var one int
	err := dbpool.QueryRow(ctx, "select 1").Scan(&one)
	if err != nil {
		return err
	}
	return nil
}

// Закрытие соединения.
func (dbpool *DBStorage) CloseConnection() {
	dbpool.Close()
}
