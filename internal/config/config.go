package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	CLI   CLIConfig
	Redis RedisConfig
	DB    DBConfig
}

type CLIConfig struct{}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR"     env-default:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB"       env-default:"0"`
}

type DBConfig struct {
	DSN          string        `env:"DB_DSN" env-default:"postgres://user:password@localhost:5432/dbname?sslmode=disable"`
	MaxOpenConns int           `env:"DB_MAX_OPEN_CONNS" env-default:"25"`
	MaxIdleConns int           `env:"DB_MAX_IDLE_CONNS" env-default:"5"`
	ConnTimeout  time.Duration `env:"DB_CONN_TIMEOUT"   env-default:"5s"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
