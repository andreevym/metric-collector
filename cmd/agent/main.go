// Package main is the entry point of the Metric Agent application.
// It initializes configurations, logging, and starts the agent.
package main

import (
	"fmt"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/transport/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	"github.com/andreevym/metric-collector/internal/config"
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

	cfg, err := config.NewAgentConfig().Init()
	if err != nil {
		log.Fatal("init agent config", err)
	}

	if _, err = logger.NewLogger(cfg.LogLevel); err != nil {
		log.Fatal("logger can't be initialized:", cfg.LogLevel, err)
	}

	// Convert poll and report intervals to time.Duration.
	pollDuration := time.Duration(cfg.PollInterval) * time.Second
	reportDuration := time.Duration(cfg.ReportInterval) * time.Second
	liveTime := time.Minute

	conn, err := grpc.Dial(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	grpcClient := proto.NewMetricCollectorClient(conn)

	// Create and run the agent.
	err = metricagent.NewAgent(
		cfg.SecretKey,
		cfg.CryptoKey,
		cfg.Address,
		pollDuration,
		reportDuration,
		liveTime,
		cfg.RateLimit,
		grpcClient,
		cfg.IsGrpcRequest,
	).Run()
	if err != nil {
		log.Fatal("failed to execute agent:", err)
	}
}
