package postgre

import "github.com/jackc/pgx"

func Connect(cnf Config) (*DBManager, error) {
	pools := make(map[string]*pgx.ConnPool)
	
	return  &DBManager{pools: pools}, nil
}

// func createUserPool(cnf userConfig) (*pgx.ConnPool, error) {
// 	conConf := pgx.ConnConfig{Host: cnf.host,
// 		Port:     cnf.port,
// 		Database: cnf.dbName,
// 		User:     cnf.user,
// 		Password: cnf.password,
// 	}
// 	conPoolConf := pgx.ConnPoolConfig{ConnConfig: conConf, MaxConnections: 15}
// 	readerPool, err := pgx.NewConnPool(conPoolConf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return readerPool, nil
// }
