package configs

import (
	"log"

	"github.com/spf13/viper"
)

const configPath = "../.env"

type DBConfig struct {
	DbName   string
	Host     string
	Port     uint16
	User     string
	Password string
}

type Config struct {
	DBConfig
}

func Load() Config {
	var cnf Config
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка считывания конфига %v", err)
	}
	dbConfig := DBConfig{
		Password: viper.GetString("POSTGRES_PASSWORD"),
		User:     viper.GetString("POSTGRES_USER"),
		Host:     "",
		Port:     viper.GetUint16(),
		DbName:   "",
	}
	cnf.DBConfig = dbConfig
	return cnf
}
