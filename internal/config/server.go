package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

type ServerConfig struct {
	Address  string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
	// StoreInterval интервал времени в секундах,
	// по истечении которого текущие показания сервера сохраняются на диск,
	// значение 0 делает запись синхронной.
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	// FileStoragePath переменная окружения FILE_STORAGE_PATH — полное имя файла,
	// куда сохраняются текущие значения,
	// пустое значение отключает функцию записи на диск.
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	// Restore определяем загружать или не загружать ранее сохранённые значения из указанного
	// файла при старте сервера.
	Restore bool `env:"RESTORE"`
	// DatabaseDsn строка с адресом подключения к БД должна получаться из переменной окружения DATABASE_DSN
	DatabaseDsn string `env:"DATABASE_DSN"`
	SecretKey   string `env:"KEY"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (c *ServerConfig) Init() *ServerConfig {
	flag.StringVar(&c.Address, "a", ":8080", "адрес и порт для запуска сервера")
	flag.StringVar(&c.LogLevel, "l", "info", "уровень логирования агента")
	flag.DurationVar(&c.StoreInterval, "i", 300*time.Second, "интервал времени в секундах "+
		"по истечении которого текущие показания сервера сохраняются на диск "+
		"(значение 0 делает запись синхронной).")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "полное имя файла "+
		"куда сохраняются текущие значения, пустое значение отключает функцию записи на диск.")
	flag.BoolVar(&c.Restore, "r", true, "определяющее, загружать или нет ранее сохранённые значения"+
		" из указанного файла при старте сервера")
	flag.StringVar(&c.DatabaseDsn, "d", "", "строка с адресом подключения к БД")
	flag.StringVar(&c.SecretKey, "k", "", "secret key, if variable is not empty will "+
		"make hash from request body and add header HashSHA256 for each http request")
	flag.Parse()

	if err := env.Parse(c); err != nil {
		panic(fmt.Errorf("failed to parse env: %w", err))
	}

	return c
}
