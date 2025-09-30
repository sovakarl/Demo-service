package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DbName   string
	Host     string
	Port     uint16
	User     string
	Password string
}

type Db struct {
	connPool *pgxpool.Pool
	log      *slog.Logger
}

func NewConnect(cnf Config, logger *slog.Logger) (*Db, error) {
	if logger == nil {
		logger = slog.Default()
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.DbName)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	
	config.MaxConns = 10
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute
	config.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	newDb := &Db{
		connPool: pool,
		log:      logger.With("component", "DB"),
	}
	
	return newDb, nil
}
