package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Netflix/go-env"
	"github.com/abdukhashimov/exporter_psql_clickhouse/pkg/logger/options"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type AppMode string

const (
	DEVELOPMENT AppMode = "DEVELOPMENT"
	PRODUCTION  AppMode = "PRODUCTION"
)

type Config struct {
	Logging options.Logging `yaml:"logging"`

	Project struct {
		Mode string `env:"APPLICATION_MODE,default=DEVELOPMENT"`
	}

	PsqlConfig struct {
		Host     string `env:"PSQL_HOST,default=localhost"`
		Port     int    `env:"PSQL_PORT,default=5432"`
		User     string `env:"PSQL_USER,default=postgres"`
		Passwrod string `env:"PSQL_PASSWORD,default=postgres"`
		Database string `env:"PSQL_DATABSE,default=sample"`
		SslMode  string `env:"PSQL_SSL_MODE,default=disable"`

		ConnString string
	}

	Network struct {
		PsqlHost string `env:"PSQL_NETWORK_HOST,default=postgres-container"`
		PsqlPort int    `env:"PSQL_NETWORK_PORT,default=5432"`

		PsqlAddress string
	}

	Exporter struct {
		TableName            string `env:"EXPORTING_TABLE_NAME"`
		ExportPerid          string `env:"EXPORTING_PERIOD,default=*/1 * * * *"`
		TelegramBotToken     string `env:"EXPORTING_BOT_TOKEN"`
		TelegramBotChannelID int64  `env:"EXPORTING_CHANNEL_ID"`
	}

	Clickhouse struct {
		Host string `env:"CLICKHOUSE_ADDRESS,default=localhost"`
		Port int    `env:"CLICKHOUSE_PORT,default=9000"`
		Auth struct {
			Database string `env:"CLICKHOUSE_DATABASE,default=default"`
			Username string `env:"CLICKHOUSE_USERNAME,default=default"`
			Password string `env:"CLICKHOUSE_PASSWORD,default="`
		}

		ConnString string
	}
}

func Load() *Config {
	var cfg Config

	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Println(".env file is not found")
	}

	appMode := getAppMode()
	configPath, err := getConfigPath(appMode)
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		panic(err)
	}

	_, err = env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		panic(err)
	}

	cfg.PsqlConfig.ConnString = cfg.makePSQLConnString()
	cfg.Clickhouse.ConnString = cfg.makeClickHouseConnString()
	cfg.Network.PsqlAddress = fmt.Sprintf("%s:%d", cfg.Network.PsqlHost, cfg.Network.PsqlPort)

	return &cfg
}

func (c *Config) makePSQLConnString() string {
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

func (c *Config) makeClickHouseConnString() string {
	return fmt.Sprintf(
		"clickhouse://%s:%d/%s?username=%s&password=%s",
		c.Clickhouse.Host,
		c.Clickhouse.Port,
		c.Clickhouse.Auth.Database,
		c.Clickhouse.Auth.Username,
		c.Clickhouse.Auth.Password,
	)
}

func getAppMode() AppMode {
	mode := AppMode(os.Getenv("APPLICATION_MODE"))
	if mode != PRODUCTION {
		mode = DEVELOPMENT
	}

	return mode
}

func getConfigPath(appMode AppMode) (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	suffix := "dev"
	if appMode == PRODUCTION {
		suffix = "prod"
	}

	return fmt.Sprintf("%s/config/config.%s.yaml", path, suffix), nil
}
