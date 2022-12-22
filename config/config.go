package config

import (
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type Config struct {
	PSQL struct {
		Host     string `env:"PSQL_HOST"`
		User     string `env:"PSQL_USER"`
		Passwrod string `env:"PSQL_PASSWORD"`
		Database string `env:"PSQL_DATABSE"`
		Table    string `env:"PSQL_TABLE"`
		Ssl      bool   `env:"PSQL_SSL"`
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
