package config

import (
	"fmt"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type Config struct {
	PsqlConfig struct {
		Host     string `env:"PSQL_HOST"`
		Port     int    `env:"PSQL_PORT"`
		User     string `env:"PSQL_USER"`
		Passwrod string `env:"PSQL_PASSWORD"`
		Database string `env:"PSQL_DATABSE"`
		SslMode  string `env:"PSQL_SSL_MODE"`
	}

	Clickhouse struct {
		Address string `env:"CLICKHOUSE_ADDRESS"`
		Auth    struct {
			Database string `env:"CLICKHOUSE_DATABASE"`
			Username string `env:"CLICKHOUSE_USERNAME"`
			Password string `env:"CLICKHOUSE_PASSWORD"`
		}
	}
}

func (c *Config) MakePSQLConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.PsqlConfig.User,
		c.PsqlConfig.Passwrod,
		c.PsqlConfig.Host,
		c.PsqlConfig.Port,
		c.PsqlConfig.Database,
		c.PsqlConfig.SslMode,
	)
}

func Load() *Config {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	_, err = env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
