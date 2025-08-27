package configs

type DBUserConfig struct {
	Role     string
	Password string
	Name     string
}

type DBConfig struct {
	DbName  string
	SSLMode string
	Host    string
	Port    uint16
	Users   []DBUserConfig
}

type Config struct {
	DBConfig
}


func Load() Config {
	return Config{}
}
