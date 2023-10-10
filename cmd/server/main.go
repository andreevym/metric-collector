package main

import (
	"flag"

	"github.com/andreevym/metric-collector/internal/server"
)

func main() {
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	server.StartServer(flagRunAddr)
}
