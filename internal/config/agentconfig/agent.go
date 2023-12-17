package agentconfig

import (
	"flag"

	"github.com/caarlos0/env"
)

type AgentConfig struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	LogLevel       string `env:"LOG_LEVEL"`
	SecretKey      string `env:"KEY"`
}

var (
	// flagAddr содержит адрес и порт для отправки метрик на сервер
	flagAddr string
	// flagSecretKey секретный ключ, если указан, то будем добавлять заголовок HashSHA256 в каждый запрос
	flagSecretKey string
	// flagReportInterval частоту отправки метрик на сервер (по умолчанию 10 секунд).
	flagReportInterval int
	// flagPollInterval частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	flagPollInterval int
	// flagLogLevel уровень логирования агента
	flagLogLevel string
)

func Flags() (*AgentConfig, error) {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagSecretKey, "k", "", "secret key, if variable is not empty will "+
		"make hash from request body and add header HashSHA256 for each http request")
	flag.IntVar(&flagReportInterval, "r", 10, "report interval (seconds)")
	flag.IntVar(&flagPollInterval, "p", 2, "poll interval (seconds)")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")

	// парсим переданные серверу агенту в зарегистрированные переменные
	flag.Parse()

	var config AgentConfig
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	if config.Address == "" {
		config.Address = flagAddr
	}

	if config.SecretKey == "" {
		config.SecretKey = flagSecretKey
	}

	// Обновлять метрики из пакета runtime с заданной частотой: pollInterval — 2 секунды.
	if config.PollInterval == 0 {
		config.PollInterval = flagPollInterval
	}

	// Отправлять метрики на сервер с заданной частотой: reportInterval — 10 секунд.
	if config.ReportInterval == 0 {
		config.ReportInterval = flagReportInterval
	}

	// Логирование, по умолчанию info
	if config.LogLevel == "" {
		config.LogLevel = flagLogLevel
	}

	return &config, nil
}
