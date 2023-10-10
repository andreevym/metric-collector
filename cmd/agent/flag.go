package main

import (
	"flag"
)

var (
	// flagAddr содержит адрес и порт для отправки метрик на сервер
	flagAddr string
	// flagReportInterval частоту отправки метрик на сервер (по умолчанию 10 секунд).
	flagReportInterval int
	// flagPollInterval частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
	flagPollInterval int
)

func init() {
	flag.StringVar(&flagAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "report interval (seconds)")
	flag.IntVar(&flagPollInterval, "p", 2, "poll interval (seconds)")
}
