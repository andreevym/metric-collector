package main

import (
	"flag"
	"time"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string

// частоту отправки метрик на сервер (по умолчанию 10 секунд).
var flagReportIntervalSec time.Duration

// частоту опроса метрик из пакета runtime (по умолчанию 2 секунды).
var flagPollIntervalSec time.Duration

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.DurationVar(&flagReportIntervalSec, "r", 10, "report interval (seconds)")
	flag.DurationVar(&flagPollIntervalSec, "p", 2, "report interval (seconds)")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
