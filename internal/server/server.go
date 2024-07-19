// Package server provides functionalities to initialize and start the metric collector server.
package server

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/andreevym/metric-collector/docs"
	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/andreevym/metric-collector/internal/storage/postgres"
	"github.com/andreevym/metric-collector/internal/storage/store"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// Start initializes and starts the metric collector server with the provided configurations.
// It sets up the metric storage based on the provided database DSN and file storage path,
// creates the necessary database client, applies migrations if required, and configures middleware.
// The server starts listening on the specified address for incoming HTTP requests.
//
// Parameters:
//   - databaseDsn: The DSN (Data Source Name) for the database connection. Leave empty for in-memory store.
//   - fileStoragePath: The path to the file storage for backup purposes. Only applicable for in-memory store.
//   - storeInterval: The interval at which metrics are stored to disk. Only applicable for in-memory store.
//   - restore: A flag indicating whether to restore metrics from the file store. Only applicable for in-memory store.
//   - secretKey: The secret key used for hashing request and response bodies.
//   - address: The address on which the HTTP server should listen (e.g., ":8080").
//
// Returns:
//   - error: An error if any occurred during initialization or while starting the HTTP server.
//
// Example:
//
//		err := Start(
//		    "postgresql://user:password@localhost:5432/database",
//		    "/path/to/file/storage",
//		    3600,
//		    true,
//		    "my-secret-key",
//	     "",
//		    ":8080",
//		)
//		if err != nil {
//		    log.Fatal(err)
//		}
func Start(
	databaseDsn string,
	fileStoragePath string,
	storeInterval int,
	restore bool,
	secretKey string,
	cryptoKey string,
	address string,
) error {
	// Initialize a background context
	ctx := context.Background()

	var metricStorage store.Storage
	var pgClient *postgres.PgClient
	var err error

	// Initialize metric storage based on the database DSN and file storage path
	if databaseDsn == "" {
		memMetricStorage := mem.NewStorage(&mem.BackupOptional{
			BackupPath:    fileStoragePath,
			StoreInterval: storeInterval,
		})

		// Restore metrics from file storage if the 'restore' flag is set
		if restore {
			err = memMetricStorage.Restore()
			if err != nil {
				return fmt.Errorf("failed to restore: %w", err)
			}
		}

		metricStorage = memMetricStorage
	} else {
		// Create a PostgreSQL client and storage
		pgClient, err = postgres.NewPgClient(databaseDsn)
		if err != nil {
			return fmt.Errorf("can't create database client: %w", err)
		}
		defer pgClient.Close()

		// Ping the database to check the connection
		err = pgClient.Ping()
		if err != nil {
			return fmt.Errorf("can't ping database: %w", err)
		}

		// Apply migrations to the database
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
			return fmt.Errorf("failed to apply migrations to the database: %w", err)
		}

		// Set metric storage to PostgreSQL storage
		metricStorage = postgres.NewPgStorage(pgClient)
	}

	// Create a new middleware instance with the provided secret key
	m := middleware.NewMiddleware(secretKey, cryptoKey)

	// Create service handlers with the initialized metric storage and PostgreSQL client
	serviceHandlers := handlers.NewServiceHandlers(metricStorage, pgClient)

	middlewares := make([]func(http.Handler) http.Handler, 0)
	if m.CryptoKey != "" {
		middlewares = append(middlewares, m.RequestCryptoMiddleware)
	}
	middlewares = append(middlewares, m.RequestGzipMiddleware)
	middlewares = append(middlewares, m.ResponseGzipMiddleware)
	middlewares = append(middlewares, m.RequestLoggerMiddleware)
	middlewares = append(middlewares, m.RequestHashMiddleware)
	middlewares = append(middlewares, m.ResponseHashMiddleware)

	// Create a new router with the service handlers and configured middleware
	router := handlers.NewRouter(serviceHandlers, middlewares...)

	// Serve Swagger UI
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))
	router.Mount("/debug", chimiddleware.Profiler())
	// Start the HTTP server and listen for incoming requests on the specified address
	return http.ListenAndServe(address, router)
}
