package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v10"
)

type AgentConfig struct {
	// Address содержит адрес и порт для отправки метрик на сервер
	Address string `env:"ADDRESS" json:"address"`
	// IsGrpcRequest если true отправляем запросы через grpc client, если false через http
	IsGrpcRequest bool `env:"IS_GRPC_REQUEST" json:"is_grpc_request"`
	// ReportInterval частоту отправки метрик на сервер
	ReportInterval int `env:"REPORT_INTERVAL" json:"report_interval"`
	// PollInterval частоту опроса метрик из пакета runtime
	PollInterval int `env:"POLL_INTERVAL" json:"poll_interval"`
	// LogLevel уровень логирования агента
	LogLevel string `env:"LOG_LEVEL" json:"log_level"`
	// SecretKey секретный ключ, если указан, то будем добавлять заголовок HashSHA256 в каждый запрос
	SecretKey string `env:"KEY" json:"key"`
	RateLimit int    `env:"RATE_LIMIT" json:"rate_limit"`
	// CryptoKey путь до файла с публичным ключом.
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{}
}

func (c *AgentConfig) GetConfigFromFile(configPath string) error {
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

func (c *AgentConfig) Init() (*AgentConfig, error) {
	flag.StringVar(&c.Address, "a", "localhost:8080", "адрес и порт для запуска сервера")
	flag.BoolVar(&c.IsGrpcRequest, "g", false, "если true отправляем запросы через grpc client, если false через http")
	flag.StringVar(&c.SecretKey, "k", "", "secret key, if variable is not empty will "+
		"make hash from request body and add header HashSHA256 for each http request")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval (seconds)")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval (seconds)")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.IntVar(&c.RateLimit, "i", 1, "количество одновременно исходящих запросов на сервер")
	flag.StringVar(&c.CryptoKey, "crypto-key", "", "путь до файла с публичным ключом")
	var configPath string
	flag.StringVar(&configPath, "config", "", "путь до конфиг файла, пример './config/agent.json'")
	flag.Parse()

	if config := os.Getenv("CONFIG"); config != "" {
		configPath = config
	}

	if configPath != "" {
		err := c.GetConfigFromFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file '%s': %w", configPath, err)
		}
	}

	if err := env.Parse(c); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return c, nil
}
