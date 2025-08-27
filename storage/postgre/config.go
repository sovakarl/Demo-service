package postgre

type Config struct {
	dbName  string
	sslMode string
	host    string
	port    uint16
	user     string
	password string
}

// type userConfig struct {
// 	role     string
// 	user     string
// 	password string
// 	Config
// }