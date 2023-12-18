package server

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andreevym/metric-collector/internal/config"
	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/andreevym/metric-collector/internal/storage/postgres"
	"go.uber.org/zap"
)

func Start(cfg *config.ServerConfig) error {
	ctx := context.Background()

	var metricStorage storage.Storage
	var pgClient *postgres.Client
	var err error

	if cfg.DatabaseDsn == "" {
		memMetricStorage := mem.NewStorage(&mem.BackupOptional{
			BackupPath:    cfg.FileStoragePath,
			StoreInterval: cfg.StoreInterval,
		})

		if cfg.Restore {
			err = memMetricStorage.Restore()
			if err != nil {
				return fmt.Errorf("failed to restore: %w", err)
			}
		}

		metricStorage = memMetricStorage
	} else {
		pgClient, err = postgres.NewClient(cfg.DatabaseDsn)
		if err != nil {
			return fmt.Errorf("can't create database client: %w", err)
		}
		defer pgClient.Close()

		err = pgClient.Ping()
		if err != nil {
			return fmt.Errorf("can't ping database: %w", err)
		}

		pgStorage := postgres.NewPgStorage(pgClient)

		err = filepath.Walk("migrations", func(path string, info fs.FileInfo, err error) error {
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
			return err
		}

		metricStorage = pgStorage
	}

	m := middleware.NewMiddleware(cfg.SecretKey)
	serviceHandlers := handlers.NewServiceHandlers(metricStorage, pgClient)
	var router = handlers.NewRouter(
		serviceHandlers,
		m.RequestGzipMiddleware,
		m.ResponseGzipMiddleware,
		m.RequestLoggerMiddleware,
		m.RequestHashMiddleware,
		m.ResponseHashMiddleware,
	)
	return http.ListenAndServe(cfg.Address, router)
}
