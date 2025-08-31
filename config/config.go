package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type DBConfig struct {
	DbName   string `mapstructure:"db_name"`
	Host     string `mapstructure:"HOST"`
	Port     uint16 `mapstructure:"port"`
	User     string `mapstructure:"pg_user"`
	Password string `mapstructure:"pg_password"`
}

type Config struct {
	DataBase DBConfig `mapstructure:",squash"`
}

func Load() (*Config, error) {
	godotenv.Load()
	
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.ReadInConfig()

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
