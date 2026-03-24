package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// Config - структура для хранения конфигурации приложения
type Config struct {
	CLI   CLIConfig
	Redis RedisConfig
}

// CLIConfig - структура для хранения конфигурации командной строки (пока пустая, но может быть расширена в будущем)
type CLIConfig struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"debug"`
}

// RedisConfig - структура для хранения конфигурации подключения к Redis
type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR"     env-default:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" env-default:""`
	DB       int    `env:"REDIS_DB"       env-default:"0"`
}

// Load - функция для загрузки конфигурации из файла .env или из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(".env", cfg); err != nil {
		if err := cleanenv.ReadEnv(cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
