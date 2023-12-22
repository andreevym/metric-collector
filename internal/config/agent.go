package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v10"
)

type AgentConfig struct {
	// Address содержит адрес и порт для отправки метрик на сервер
	Address string `env:"ADDRESS"`
	// ReportInterval частоту отправки метрик на сервер
	ReportInterval int `env:"REPORT_INTERVAL"`
	// PollInterval частоту опроса метрик из пакета runtime
	PollInterval int `env:"POLL_INTERVAL"`
	// LogLevel уровень логирования агента
	LogLevel string `env:"LOG_LEVEL"`
	// SecretKey секретный ключ, если указан, то будем добавлять заголовок HashSHA256 в каждый запрос
	SecretKey string `env:"KEY"`
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{}
}

func (c *AgentConfig) Init() *AgentConfig {
	flag.StringVar(&c.Address, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&c.SecretKey, "k", "", "secret key, if variable is not empty will "+
		"make hash from request body and add header HashSHA256 for each http request")
	flag.IntVar(&c.ReportInterval, "r", 10, "report interval (seconds)")
	flag.IntVar(&c.PollInterval, "p", 2, "poll interval (seconds)")
	flag.StringVar(&c.LogLevel, "l", "info", "log level")
	flag.Parse()

	if err := env.Parse(c); err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	return c
}
