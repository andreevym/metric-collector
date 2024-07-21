package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
)

type ServerConfig struct {
	// Address адрес и порт для запуска сервера
	Address string `env:"ADDRESS" json:"address"`
	// LogLevel уровень логирования
	LogLevel string `env:"LOG_LEVEL" json:"log_level"`
	// StoreInterval интервал времени в секундах,
	// по истечении которого текущие показания сервера сохраняются на диск,
	// значение 0 делает запись синхронной.
	StoreInterval int `env:"STORE_INTERVAL" json:"store_interval"`
	// FileStoragePath переменная окружения FILE_STORAGE_PATH — полное имя файла,
	// куда сохраняются текущие значения,
	// пустое значение отключает функцию записи на диск.
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	// Restore определяем загружать или не загружать ранее сохранённые значения из указанного
	// файла при старте сервера.
	Restore bool `env:"RESTORE" json:"restore"`
	// DatabaseDsn строка с адресом подключения к БД должна получаться из переменной окружения DATABASE_DSN
	DatabaseDsn string `env:"DATABASE_DSN" json:"database_dsn"`
	// SecretKey секретный ключ, если переменная не пустая "+
	// тогда добавляем в заголовок каждого запроса hash от request body под ключом HashSHA256
	SecretKey string `env:"KEY" json:"key"`
	// CryptoKey путь до файла с публичным ключом.
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
	// TrustedSubnet строковое представление бесклассовой адресации (CIDR)
	TrustedSubnet string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (c *ServerConfig) GetConfigFromFile(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read file by path '%s': %w", configPath, err)
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data '%s': %w", string(data), err)
	}

	return nil
}

func (c *ServerConfig) Init() *ServerConfig {
	flag.StringVar(&c.Address, "a", ":8080", "адрес и порт для запуска сервера")
	flag.StringVar(&c.LogLevel, "l", "info", "уровень логирования")
	flag.IntVar(&c.StoreInterval, "i", 300, "интервал времени в секундах "+
		"по истечении которого текущие показания сервера сохраняются на диск "+
		"(значение 0 делает запись синхронной).")
	flag.StringVar(&c.FileStoragePath, "f", "/tmp/metrics-db.json", "полное имя файла "+
		"куда сохраняются текущие значения, пустое значение отключает функцию записи на диск.")
	flag.BoolVar(&c.Restore, "r", true, "определяющее, загружать или нет ранее сохранённые значения"+
		" из указанного файла при старте сервера")
	flag.StringVar(&c.DatabaseDsn, "d", "", "строка с адресом подключения к БД")
	flag.StringVar(&c.SecretKey, "k", "", "секретный ключ, если переменная не пустая "+
		"тогда добавляем в заголовок каждого запроса hash от request body под ключом HashSHA256")
	flag.StringVar(&c.CryptoKey, "crypto-key", "", "путь до файла с публичным ключом")
	flag.StringVar(&c.TrustedSubnet, "t", "", "строковое представление бесклассовой адресации (CIDR)")
	var configPath string
	flag.StringVar(&configPath, "config", "", "путь до конфиг файла, пример './config/server.json'")
	flag.Parse()

	if config := os.Getenv("CONFIG"); config != "" {
		configPath = config
	}

	if configPath != "" {
		err := c.GetConfigFromFile(configPath)
		if err != nil {
			panic(fmt.Errorf("failed to read config file '%s': %w", configPath, err))
		}
	}

	if err := env.Parse(c); err != nil {
		panic(fmt.Errorf("failed to parse env: %w", err))
	}

	return c
}
