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
		counterStorage = mem.NewStorage()
		gaugeStorage = mem.NewStorage()
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

		err := filepath.Walk("migrations", func(path string, info fs.FileInfo, err error) error {
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

	store, err := multistorage.NewMetricStorage(counterStorage, gaugeStorage, cfg)
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
