package main

import (
	"flag"
	"time"
)

// flagAddr содержит адрес и порт для отправки метрик на сервер
var flagAddr string

// flagReportInterval частоту отправки метрик на сервер (по умолчанию 10 секунд).
var flagReportInterval time.Duration

// flagPollInterval частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
var flagPollInterval time.Duration

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	flag.StringVar(&flagAddr, "a", "http://localhost:8080", "address and port to run server")
	flag.DurationVar(&flagReportInterval, "r", 10, "report interval (seconds)")
	flag.DurationVar(&flagPollInterval, "p", 2, "report interval (seconds)")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
