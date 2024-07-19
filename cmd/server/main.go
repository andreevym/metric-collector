// Package main is the entry point of the Metric Collector application.
// It initializes configurations, logging, and starts the server.
package main

import (
	"fmt"
	"log"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/server"
)

var buildVersion string
var buildDate string
var buildCommit string

func printVersion() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

// @title Metric Collector API
// @version 18.0
// @description Metrics and Alerting Service
// @termsOfService http://swagger.io/terms/
// @contact.name Metric Collector API Support
// @contact.url http://www.swagger.io/support
// @contact.email andreevym@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// main is the entry point of the application.
// It initializes configurations, logging, and starts the server.
func main() {
	printVersion()

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
		cfg.CryptoKey,
		cfg.Address,
	)
	if err != nil {
		logger.Logger().Fatal(err.Error())
	}
}
