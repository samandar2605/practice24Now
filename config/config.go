package config

import (
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type Config struct {
	PostConfig PostgresConfig
	RedisConfig RedisConfig
	HttpPort   string
	SMTP Smtp
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Database string
	Password string
}

type RedisConfig struct{
	RedisHost           string
	RedisPort           string
}

type Smtp struct{
	Sender string
	Password string
}


func Load(path string) Config {
	gotenv.Load(".env")
	Conf := viper.New()
	Conf.AutomaticEnv()
	cfg := Config{
		HttpPort: Conf.GetString("HTTP_PORT"),
		PostConfig: PostgresConfig{
			Host:     Conf.GetString("POSTGRES_HOST"),
			Port:     Conf.GetString("POSTGRES_PORT"),
			User:     Conf.GetString("POSTGRES_USER"),
			Database: Conf.GetString("POSTGRES_DATABASE"),
			Password: Conf.GetString("POSTGRES_PASSWORD"),
		},
		RedisConfig: RedisConfig{
			RedisHost: Conf.GetString("REDIS_HOST"),
			RedisPort: Conf.GetString("REDIS_PORT"),
		},
		SMTP: Smtp{
			Sender: Conf.GetString("SMTP_SENDER"),
			Password: Conf.GetString("SMTP_PASSWORD"),
		},
	}
	return cfg
}
