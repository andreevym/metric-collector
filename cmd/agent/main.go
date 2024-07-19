// Package main is the entry point of the Metric Agent application.
// It initializes configurations, logging, and starts the agent.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/metricagent"
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

// main is the entry point of the application.
// It initializes configurations, logging, and starts the agent.
func main() {
	printVersion()

	// Initialize agent configurations.
	cfg := config.NewAgentConfig().Init()
	if cfg == nil {
		log.Fatal("agent config can't be nil")
	}

	// Initialize logger.
	_, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal("logger can't be initialized:", cfg.LogLevel, err)
	}

	// Convert poll and report intervals to time.Duration.
	pollDuration := time.Duration(cfg.PollInterval) * time.Second
	reportDuration := time.Duration(cfg.ReportInterval) * time.Second
	liveTime := time.Minute

	// Create and run the agent.
	err = metricagent.NewAgent(
		cfg.SecretKey,
		cfg.CryptoKey,
		cfg.Address,
		pollDuration,
		reportDuration,
		liveTime,
		cfg.RateLimit,
	).Run()
	if err != nil {
		log.Fatal("failed to execute agent:", err)
	}
}
