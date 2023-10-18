package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

// flagRunAddr адрес и порт для запуска сервера, аргумент -a со значением :8080 по умолчанию
var flagRunAddr string

func ServerFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
}

func ServerParse() (*ServerConfig, error) {
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

	return &config, nil
}
