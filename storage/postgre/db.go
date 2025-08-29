package postgre

import "github.com/jackc/pgx"


func Connect(cnf Config) (*pgx.ConnPool, error) {
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

	return pool, nil
}

