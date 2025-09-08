package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type DB struct {
	DbName   string `mapstructure:"db_name"`
	Host     string `mapstructure:"db_host"`
	Port     uint16 `mapstructure:"db_port"`
	User     string `mapstructure:"pg_user"`
	Password string `mapstructure:"pg_password"`
}
type App struct {
	Host             string `mapstructure:"app_host"`
	Port             uint16 `mapstructure:"app_port"`
	CacheWarmUpLimit uint64 `mapstructure:"app_cacheWarmUpLimit"`
}

// type BrokerConfig struct {
// }

type Logger struct {
	LogLevel string `mapstructure:"log_lvl"`
}

type Config struct {
	DataBase DB     `mapstructure:",squash"`
	App      App    `mapstructure:",squash"`
	Log      Logger `mapstructure:",squash"`
}

// func (c *Config) GetLogConfig() Logger {
// 	return c.Log
// }

// func (c *Config) GetAppConfig() App {
// 	return c.App
// }

// func (c *Config) GetDbConfig() DB {
// 	return c.DataBase
// }

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
