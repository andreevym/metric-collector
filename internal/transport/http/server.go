// Package transport provides functionalities to initialize and start the metric collector transport.
package http

import (
	"context"
	"fmt"
	"github.com/andreevym/metric-collector/internal/logger"
	"go.uber.org/zap"
	"net"
	"net/http"

	_ "github.com/andreevym/metric-collector/docs"
	"github.com/andreevym/metric-collector/internal/storage/postgres"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"github.com/andreevym/metric-collector/internal/transport/http/handlers"
	"github.com/andreevym/metric-collector/internal/transport/http/middleware"
)

type Server struct {
	handler http.Handler
	server  *http.Server
	address string
}

func NewHTTPServer(pgClient *postgres.PgClient, metricStorage store.Storage, secretKey string, cryptoKey string, trustedSubnet string, address string) (*Server, error) {
	var err error
	var ipTrustedSubnet *net.IPNet
	if trustedSubnet != "" {
		_, ipTrustedSubnet, err = net.ParseCIDR(trustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("failed to parse trusted subnet: %w", err)
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
	return &Server{handler: router, address: address}, nil
}

func (s *Server) Run() error {
	s.server = &http.Server{Addr: s.address, Handler: s.handler}
	logger.Logger().Info("listening http server", zap.String("address", s.address))
	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Logger().Info("shutting down http server")
	return s.server.Shutdown(ctx)
}
