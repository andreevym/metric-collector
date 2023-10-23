package serverconfig

import (
	"flag"

	"github.com/caarlos0/env"
)

type ServerConfig struct {
	Address  string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
}

var (
	// flagRunAddr адрес и порт для запуска сервера, аргумент -a со значением :8080 по умолчанию
	flagRunAddr string
	// flagLogLevel уровень логирования агента
	flagLogLevel string
)

func Flags() (*ServerConfig, error) {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	var config ServerConfig

	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	if config.Address == "" {
		config.Address = flagRunAddr
	}

	// Логирование, по умолчанию info
	if config.LogLevel == "" {
		config.LogLevel = flagLogLevel
	}

	return &config, nil
}
