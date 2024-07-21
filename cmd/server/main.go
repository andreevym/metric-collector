// Package main is the entry point of the Metric Collector application.
// It initializes configurations, logging, and starts the server.
package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/server"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/andreevym/metric-collector/internal/storage/postgres"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"go.uber.org/zap"
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

	ctx, cancel := context.WithCancel(context.Background())

	pgClient, err := BuildPgClient(ctx, cfg)
	if err != nil {
		logger.Logger().Fatal("can't build PgClient", zap.Error(err))
	}
	defer pgClient.Close()

	storage, err := BuildStorage(pgClient, cfg)
	if err != nil {
		logger.Logger().Fatal("can't create metric storage", zap.Error(err))
	}

	s := server.NewServer(pgClient, storage, cfg.SecretKey, cfg.CryptoKey, cfg.TrustedSubnet)
	go func() {
		defer cancel()
		s.Run(cfg.Address)
	}()

	s.WaitShutdown(ctx, storage)
}

func BuildPgClient(ctx context.Context, cfg *config.ServerConfig) (*postgres.PgClient, error) {
	if cfg.DatabaseDsn == "" {
		return nil, nil
	}
	// Create a PostgreSQL client and storage
	pgClient, err := postgres.NewPgClient(cfg.DatabaseDsn)
	if err != nil {
		return nil, fmt.Errorf("can't create database client: %w", err)
	}
	// Ping the database to check the connection
	err = pgClient.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping database: %w", err)
	}

	err = RunMigrationPostgres(ctx, pgClient)
	if err != nil {
		return nil, fmt.Errorf("can't migrate database: %w", err)
	}
	return pgClient, nil
}

func BuildStorage(pgClient *postgres.PgClient, cfg *config.ServerConfig) (store.Storage, error) {
	if pgClient != nil {
		return postgres.NewPgStorage(pgClient), nil
	}

	memMetricStorage := mem.NewStorage(&mem.BackupOptional{
		BackupPath:    cfg.FileStoragePath,
		StoreInterval: cfg.StoreInterval,
	})

	// Restore metrics from file storage if the 'restore' flag is set
	if cfg.Restore {
		err := memMetricStorage.Restore()
		if err != nil {
			return nil, fmt.Errorf("failed to restore: %w", err)
		}
	}

	return memMetricStorage, nil
}

func RunMigrationPostgres(ctx context.Context, pgClient *postgres.PgClient) error {
	// Apply migrations to the database
	err := filepath.Walk("migrations", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			logger.Logger().Info("apply migration", zap.String("path", path))
			bytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			err = pgClient.ApplyMigration(ctx, string(bytes))
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to apply migrations to the database: %w", err)
	}

	return nil
}
