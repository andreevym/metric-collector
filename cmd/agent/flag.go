package main

import (
	"flag"
	"time"
)

var (
	// flagAddr содержит адрес и порт для отправки метрик на сервер
	flagAddr string
	// flagReportInterval частоту отправки метрик на сервер (по умолчанию 10 секунд).
	flagReportInterval time.Duration
	// flagPollInterval частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	flagPollInterval time.Duration
)

func init() {
	flag.StringVar(&flagAddr, "a", "http://localhost:8080", "address and port to run server")
	flag.DurationVar(&flagReportInterval, "r", 10, "report interval (seconds)")
	flag.DurationVar(&flagPollInterval, "p", 2, "poll interval (seconds)")
}
