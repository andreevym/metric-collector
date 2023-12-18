package server

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/andreevym/metric-collector/internal/storage/postgres"
	"go.uber.org/zap"
)

func Start(
	databaseDsn string,
	fileStoragePath string,
	storeInterval int,
	restore bool,
	secretKey string,
	address string,
) error {
	ctx := context.Background()

	var metricStorage storage.Storage
	var pgClient *postgres.Client
	var err error

	if databaseDsn == "" {
		metricStorage, err = buildMemStorage(fileStoragePath, storeInterval, restore)
		if err != nil {
			return fmt.Errorf("failed to build mem storage: %w", err)
		}
	} else {
		pgClient, metricStorage, err = buildPostgresStorage(ctx, databaseDsn)
		if err != nil {
			return fmt.Errorf("failed to build postgres storage: %w", err)
		}
	}

	m := middleware.NewMiddleware(secretKey)
	serviceHandlers := handlers.NewServiceHandlers(metricStorage, pgClient)
	var router = handlers.NewRouter(
		serviceHandlers,
		m.RequestGzipMiddleware,
		m.ResponseGzipMiddleware,
		m.RequestLoggerMiddleware,
		m.RequestHashMiddleware,
		m.ResponseHashMiddleware,
	)
	return http.ListenAndServe(address, router)
}

func buildMemStorage(fileStoragePath string, storeInterval int, restore bool) (*mem.Storage, error) {
	memMetricStorage := mem.NewStorage(&mem.BackupOptional{
		BackupPath:    fileStoragePath,
		StoreInterval: storeInterval,
	})

	if restore {
		err := memMetricStorage.Restore()
		if err != nil {
			return nil, fmt.Errorf("failed to restore: %w", err)
		}
	}
	return memMetricStorage, nil
}

func buildPostgresStorage(
	ctx context.Context,
	databaseDsn string,
) (*postgres.Client, *postgres.PgStorage, error) {
	pgClient, err := postgres.NewClient(databaseDsn)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create database client: %w", err)
	}
	defer pgClient.Close()

	err = pgClient.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("can't ping database: %w", err)
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
		return nil, nil, err
	}
	return pgClient, pgStorage, nil
}
