package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
)

type ServerConfig struct {
	Address     string        `env:"RUN_ADDRESS" envDefault:"0.0.0.0:8080"`
	DatabaseDSN string        `env:"DATABASE_URI,unset" envDefault:"host=postgres-gophermart port=25432 user=username password=password dbname=gophermart sslmode=disable"`
	LogLevel    string        `env:"LOG_LEVEL" envDefault:"debug"`
	TokenTTL    time.Duration `env:"TOKEN_TTL" envDefault:"1h"`
	SecretKey   string        `env:"SECRET_KEY,unset" envDefault:"local-default-secret"`
}

type AccrualConfig struct {
	Address      string        `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://accrual:8090"`
	Interval     time.Duration `env:"ACCRUAL_PROCESSING_INTERVAL" envDefault:"5s"`
	WorkersCount int           `env:"WORKERS_COUNT" envDefault:"3"`
}

type Config struct {
	Server  ServerConfig
	Accrual AccrualConfig
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	var serverAddress, databaseDSN, logLevel, address string
	flag.StringVar(&serverAddress, "a", "", "Address and port to run server")
	flag.StringVar(&databaseDSN, "d", "", "Database URI")
	flag.StringVar(&logLevel, "l", "", "Log levle: debug, info, warn, error, panic, fatal")
	flag.StringVar(&address, "r", "", "Accrual system addres")
	flag.Parse()

	if serverAddress != "" {
		cfg.Server.Address = serverAddress
	}
	if databaseDSN != "" {
		cfg.Server.DatabaseDSN = databaseDSN
	}
	if logLevel != "" {
		cfg.Server.LogLevel = logLevel
	}
	if address != "" {
		cfg.Accrual.Address = address
	}

	return &cfg, nil
}
