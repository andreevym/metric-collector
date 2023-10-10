package main

import (
	"flag"
	"log"

	"github.com/andreevym/metric-collector/internal/server"
	"github.com/caarlos0/env"
)

func main() {
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	var config EnvConfig
	err := env.Parse(&config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Address != "" {
		flagRunAddr = config.Address
	}

	server.StartServer(flagRunAddr)
}
