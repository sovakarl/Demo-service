package postgres

import "github.com/jackc/pgx"

type Db struct {
	connPool *pgx.ConnPool
}

func (db Db) Close() {
	db.connPool.Close()
}

func NewConnect(cnf Config) (*Db, error) {
	conConfig := pgx.ConnConfig{Host: cnf.Host,
		Port:     cnf.Port,
		Database: cnf.DbName,
		Password: cnf.Password,
		User:     cnf.User,
	}

	poolConfig := pgx.ConnPoolConfig{ConnConfig: conConfig}

	pool, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		return nil, err
	}

	return &Db{connPool: pool}, nil
}
