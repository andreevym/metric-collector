// Package main is the entry point of the Metric Collector application.
// It initializes configurations, logging, and starts the server.
package main

import (
	"log"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/server"
)

// main is the entry point of the application.
// It initializes configurations, logging, and starts the server.
func main() {
	// Initialize server configurations.
	cfg := config.NewServerConfig().Init()
	if cfg == nil {
		log.Fatal("server config can't be nil")
	}

	// Initialize logger.
	_, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal("logger can't be initialized:", cfg.LogLevel, err)
	}

	// Start the server.
	err = server.Start(
		cfg.DatabaseDsn,
		cfg.FileStoragePath,
		cfg.StoreInterval,
		cfg.Restore,
		cfg.SecretKey,
		cfg.Address,
	)
	if err != nil {
		logger.Logger().Fatal(err.Error())
	}
}
