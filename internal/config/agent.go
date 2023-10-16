package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type AgentConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

var (
	// flagAddr содержит адрес и порт для отправки метрик на сервер
	flagAddr string
	// flagReportInterval частоту отправки метрик на сервер (по умолчанию 10 секунд).
	flagReportInterval int
	// flagPollInterval частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	flagPollInterval int
)

func AgentParse() (*AgentConfig, error) {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "report interval (seconds)")
	flag.IntVar(&flagPollInterval, "p", 2, "poll interval (seconds)")

	// парсим переданные серверу агенту в зарегистрированные переменные
	flag.Parse()

	var config *AgentConfig
	err := env.Parse(config)
	if err != nil {
		return nil, err
	}

	if config.Address == "" {
		config.Address = flagAddr
	}

	// Обновлять метрики из пакета runtime с заданной частотой: pollInterval — 2 секунды.
	if config.PollInterval == 0 {
		config.PollInterval = flagPollInterval
	}

	// Отправлять метрики на сервер с заданной частотой: reportInterval — 10 секунд.
	if config.ReportInterval == 0 {
		config.ReportInterval = flagReportInterval
	}

	return config, nil
}
