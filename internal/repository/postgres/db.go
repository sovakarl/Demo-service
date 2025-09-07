package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	connPool *pgxpool.Pool
}

func (db Db) Close() {
	db.connPool.Close()
}

func NewConnect(cnf Config) (*Db, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.DbName)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return &Db{connPool: pool}, nil
}
