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
	Addres       string        `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://accrual:8090"`
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

	if value := flag.String("a", "", "Address and port to run server"); *value != "" {
		cfg.Server.Address = *value
	}
	if value := flag.String("d", "", "Database URI"); *value != "" {
		cfg.Server.DatabaseDSN = *value
	}
	if value := flag.String("l", "", "Log levle: debug, info, warn, error, panic, fatal"); *value != "" {
		cfg.Server.LogLevel = *value
	}
	if value := flag.String("r", "", "Accrual system addres"); *value != "" {
		cfg.Accrual.Addres = *value
	}
	// flag.Parse()  // todo ????

	return &cfg, nil
}
