package config

import (
	"github.com/spf13/viper"
)

type DB struct {
	DbName   string `mapstructure:"DB_NAME"`
	Host     string `mapstructure:"DB_HOST"`
	Port     uint16 `mapstructure:"DB_PORT"`
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
	DataBase DB     `mapstructure:",squash"`
	App      App    `mapstructure:",squash"`
	Log      Logger `mapstructure:",squash"`
}

func Load() (*Config, error) {

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")

	viper.BindEnv("APP_HOST")
	viper.BindEnv("APP_PORT")

	viper.ReadInConfig()

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
