package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andreevym/metric-collector/internal/config/serverconfig"
	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/pg"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/andreevym/metric-collector/internal/storage/postgres"
	"go.uber.org/zap"
)

func Start(cfg *serverconfig.ServerConfig) error {
	var counterStorage storage.Storage
	var gaugeStorage storage.Storage
	var dbClient *pg.Client
	var err error

	if cfg.DatabaseDsn == "" {

		var counterBackupPath string
		var gaugeBackupPath string

		if cfg.FileStoragePath != "" {
			err = os.MkdirAll(cfg.FileStoragePath+"/", 0777)
			if err != nil {
				return err
			}
			if ok, _ := isDirectory(cfg.FileStoragePath); !ok {
				return fmt.Errorf("storage path need to be directory %s", cfg.FileStoragePath)
			}
			counterBackupPath = cfg.FileStoragePath + "/counter.backup"
			gaugeBackupPath = cfg.FileStoragePath + "/gauge.backup"
		}
		memCounterStorage := mem.NewStorage(&mem.BackupOptional{
			BackupPath:    counterBackupPath,
			StoreInterval: cfg.StoreInterval,
		})
		memGaugeStorage := mem.NewStorage(&mem.BackupOptional{
			BackupPath:    gaugeBackupPath,
			StoreInterval: cfg.StoreInterval,
		})

		if cfg.Restore {
			err = memCounterStorage.Restore()
			if err != nil {
				return fmt.Errorf("failed to restore: %w", err)
			}
			err = memGaugeStorage.Restore()
			if err != nil {
				return fmt.Errorf("failed to restore: %w", err)
			}
		}

		counterStorage = memCounterStorage
		gaugeStorage = memGaugeStorage
	} else {
		dbClient, err = pg.NewClient(cfg.DatabaseDsn)
		if err != nil {
			return fmt.Errorf("can't create database client: %w", err)
		}
		defer dbClient.Close()

		err = dbClient.Ping()
		if err != nil {
			return fmt.Errorf("can't ping database: %w", err)
		}

		counterPgStorage := postgres.NewPgStorage(dbClient, "counter")
		gaugePgStorage := postgres.NewPgStorage(dbClient, "gauge")

		err = filepath.Walk("migrations", func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				logger.Log.Info("apply migration", zap.String("path", path))
				bytes, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				err = dbClient.ApplyMigration(string(bytes))
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}

		counterStorage = counterPgStorage
		gaugeStorage = gaugePgStorage
	}

	store, err := multistorage.NewMetricManager(counterStorage, gaugeStorage)
	if err != nil {
		return err
	}

	serviceHandlers := handlers.NewServiceHandlers(store, dbClient)
	router := handlers.NewRouter(
		serviceHandlers,
		middleware.GzipRequestMiddleware,
		middleware.GzipResponseMiddleware,
		middleware.RequestLogger,
	)

	return http.ListenAndServe(cfg.Address, router)
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
