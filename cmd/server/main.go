// Package main is the entry point of the Metric Collector application.
// It initializes configurations, logging, and starts the server.
package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

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

	cfg, err := config.NewServerConfig().Init()
	if err != nil {
		log.Fatal("init server config", err)
	}

	if _, err := logger.NewLogger(cfg.LogLevel); err != nil {
		log.Fatal("logger can't be initialized:", cfg.LogLevel, err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	pgClient, err := BuildPgClient(ctx, cfg)
	if err != nil {
		logger.Logger().Fatal("can't build PgClient", zap.Error(err))
	}
	defer pgClient.Close()

	storeInterval := time.Duration(cfg.StoreInterval) * time.Second
	storage, err := BuildStorage(pgClient, cfg, storeInterval)
	if err != nil {
		logger.Logger().Fatal("can't create metric storage", zap.Error(err))
	}

	s := server.NewServer(pgClient, storage, cfg.SecretKey, cfg.CryptoKey, cfg.TrustedSubnet)
	go func() {
		defer cancel()
		s.Run(cfg.Address)
	}()

	if storeInterval > 0 {
		go func() {
			for {
				ticker := time.NewTicker(storeInterval)
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := storage.BackupPeriodically(); err != nil {
						logger.Logger().Fatal("Backup failed", zap.Error(err))
					}
				}
			}
		}()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	for {
		select {
		case <-quit:
			fmt.Println("Shutting down server...")

			ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()

			if err := s.Server.Shutdown(ctx); err != nil {
				log.Fatalf("Server shutdown failed: %v", err)
			}
			fmt.Println("Server stopped gracefully")

			if err := storage.Backup(); err != nil {
				logger.Logger().Fatal("Backup failed", zap.Error(err))
			}
			return
		case <-ctx.Done():
			logger.Logger().Info("Shutting down server...")
			return
		}
	}
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
	if err = pgClient.Ping(); err != nil {
		return nil, fmt.Errorf("can't ping database: %w", err)
	}

	if err = applyMigrations(ctx, pgClient); err != nil {
		return nil, fmt.Errorf("failed to apply migrations to the database: %w", err)
	}

	return pgClient, nil
}

func BuildStorage(pgClient *postgres.PgClient, cfg *config.ServerConfig, storeInterval time.Duration) (store.Storage, error) {
	if pgClient != nil {
		return postgres.NewPgStorage(pgClient), nil
	}

	memMetricStorage := mem.NewStorage(&mem.BackupOptional{
		BackupPath:    cfg.FileStoragePath,
		StoreInterval: storeInterval,
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

func applyMigrations(ctx context.Context, pgClient *postgres.PgClient) error {
	return filepath.Walk("migrations", func(path string, info fs.FileInfo, err error) error {
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
}
