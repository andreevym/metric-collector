package main

import (
	"log"
	"time"

	"github.com/andreevym/metric-collector/internal/config/agentconfig"
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

	pollDuration := time.Duration(cfg.PollInterval) * time.Second
	reportDuration := time.Duration(cfg.ReportInterval) * time.Second

	metric.Start(pollDuration, reportDuration, cfg.Address)
}
