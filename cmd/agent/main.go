package main

import (
	"log"
	"time"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/metricagent"
)

func main() {
	cfg := config.NewAgentConfig().Init()
	if cfg == nil {
		log.Fatal("agent config can't be nil")
	}
	_, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal("logger can't be init", cfg.LogLevel, err)
	}

	pollDuration := time.Duration(cfg.PollInterval) * time.Second
	reportDuration := time.Duration(cfg.ReportInterval) * time.Second
	liveTime := time.Minute

	a := metricagent.NewAgent(cfg.SecretKey, cfg.Address, pollDuration, reportDuration, liveTime)
	a.Start()
}
