package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DB struct {
	DbName   string `mapstructure:"DB_NAME"`
	Host     string `mapstructure:"APP_DB_HOST"`
	Port     uint16 `mapstructure:"APP_DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
}
type App struct {
	Host string `mapstructure:"APP_HOST"`
	Port uint16 `mapstructure:"APP_PORT"`
}

type Logger struct {
	LogLevel string `mapstructure:"LOG_LVL"`
}

type Config struct {
	DataBase DB       `mapstructure:",squash"`
	App      App      `mapstructure:",squash"`
	Log      Logger   `mapstructure:",squash"`
	Consumer Consumer `mapstructure:",squash"`
}

type Consumer struct {
	Topic   string `mapstructure:"KAFKA_TOPIC"`
	Host    string `mapstructure:"KAFKA_HOST"`
	Port    uint16 `mapstructure:"KAFKA_PORT"`
	GroupID string `mapstructure:"KAFKA_GroupID"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	viper.BindEnv("APP_DB_PORT")
	viper.BindEnv("APP_DB_HOST")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")

	viper.BindEnv("KAFKA_TOPIC")
	viper.BindEnv("KAFKA_HOST")
	viper.BindEnv("KAFKA_PORT")
	viper.BindEnv("KAFKA_GroupID")

	viper.BindEnv("APP_HOST")
	viper.BindEnv("APP_PORT")
	viper.BindEnv("LOG_LVL")

	viper.ReadInConfig()

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	fmt.Println(cfg)
	return &cfg, nil
}
