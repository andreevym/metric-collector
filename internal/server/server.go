// Package server provides functionalities to initialize and start the metric collector server.
package server

import (
	"fmt"
	"net"
	"net/http"

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
	trustedSubnet string,
) *Server {
	var err error
	var ipTrustedSubnet *net.IPNet
	if trustedSubnet != "" {
		_, ipTrustedSubnet, err = net.ParseCIDR(trustedSubnet)
		if err != nil {
			panic(fmt.Errorf("failed to parse trusted subnet: %w", err))
		}
	}
	m := middleware.NewMiddleware(secretKey, cryptoKey, ipTrustedSubnet)
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
	if m.TrustedSubnet != nil {
		middlewares = append(middlewares, m.TrustedSubnetMiddleware)
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
