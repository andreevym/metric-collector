package main

import (
	"log"

	"github.com/andreevym/metric-collector/internal/config/serverconfig"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/server"
)

func main() {
	cfg, err := serverconfig.Flags()
	if err != nil {
		log.Fatal(err)
	}
	if cfg == nil {
		log.Fatal("server config can't be nil")
	}

	_, err = logger.Logger(cfg.LogLevel)
	if err != nil {
		log.Fatal("logger can't be init", cfg.LogLevel, err)
	}

	err = server.Start(cfg.Address)
	if err != nil {
		log.Fatal(err)
	}
}
