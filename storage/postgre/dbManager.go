package postgre

import "github.com/jackc/pgx"

type DBManager struct {
	pools map[string]*pgx.ConnPool
}

// func NewDBManager()

func (db *DBManager) Read(id int) {
	//
}

func (db *DBManager) Write() {
	//
}

func (db *DBManager) Close() {
	for _, pool := range db.pools {
		pool.Close()
	}
}
