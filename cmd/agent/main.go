// Package main is the entry point of the Metric Agent application.
// It initializes configurations, logging, and starts the agent.
package main

import (
	"log"
	"time"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/metricagent"
)

// main is the entry point of the application.
// It initializes configurations, logging, and starts the agent.
func main() {
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
