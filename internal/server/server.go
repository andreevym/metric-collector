// Package server provides functionalities to initialize and start the metric collector server.
package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/andreevym/metric-collector/docs"
	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/storage/postgres"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"go.uber.org/zap"
)

type Server struct {
	Handler http.Handler
	Server  *http.Server
}

func NewServer(
	pgClient *postgres.PgClient,
	metricStorage store.Storage,
	secretKey string,
	cryptoKey string,
) *Server {
	m := middleware.NewMiddleware(secretKey, cryptoKey)
	serviceHandlers := handlers.NewServiceHandlers(metricStorage, pgClient)
	middlewares := []func(http.Handler) http.Handler{
		m.RequestGzipMiddleware,
		m.ResponseGzipMiddleware,
		m.RequestLoggerMiddleware,
		m.RequestHashMiddleware,
		m.ResponseHashMiddleware,
	}
	if m.CryptoKey != "" {
		middlewares = append(middlewares, m.RequestCryptoMiddleware)
	}
	router := handlers.NewRouter(serviceHandlers, middlewares...)
	return &Server{Handler: router}
}

func (s *Server) Run(addr string) {
	s.Server = &http.Server{Addr: addr, Handler: s.Handler}
	if err := s.Server.ListenAndServe(); err != nil {
		logger.Logger().Fatal(fmt.Sprintf("failed to start server: %v", err))
	}
	logger.Logger().Info("Server listening", zap.String("addr", addr))
}

func (s *Server) WaitShutdown(ctx context.Context, storage store.Storage) {
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
			err := storage.Backup()
			if err != nil {
				logger.Logger().Fatal("Backup failed", zap.Error(err))
			}
		case <-ctx.Done():
			logger.Logger().Info("Shutting down server...")
			return
		}
	}
}
