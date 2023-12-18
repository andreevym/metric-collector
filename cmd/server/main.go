package main

import (
	"log"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/server"
)

func main() {
	cfg := config.NewServerConfig().Init()
	if cfg == nil {
		log.Fatal("server config can't be nil")
	}

	_, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal("logger can't be init", cfg.LogLevel, err)
	}

	err = server.Start(cfg)
	if err != nil {
		logger.Logger().Fatal(err.Error())
	}
}
