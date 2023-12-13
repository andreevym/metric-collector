package main

import (
	"log"
	"time"

	"github.com/andreevym/metric-collector/internal/config/agentconfig"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/metric"
)

func main() {
	cfg, err := agentconfig.Flags()
	if err != nil {
		log.Fatal(err)
	}
	if cfg == nil {
		log.Fatal("agent config can't be nil")
	}
	_, err = logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal("logger can't be init", cfg.LogLevel, err)
	}

	pollDuration := time.Duration(cfg.PollInterval) * time.Second
	reportDuration := time.Duration(cfg.ReportInterval) * time.Second

	metric.Start(pollDuration, reportDuration, cfg.Address)
}
