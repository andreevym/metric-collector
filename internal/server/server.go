package server

import (
	"net/http"

	"github.com/andreevym/metric-collector/internal/config/serverconfig"
	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
)

func Start(cfg *serverconfig.ServerConfig) error {
	counterMemStorage := mem.NewStorage()
	gaugeMemStorage := mem.NewStorage()
	store, err := multistorage.NewStorage(counterMemStorage, gaugeMemStorage, cfg)
	if err != nil {
		return err
	}
	serviceHandlers := handlers.NewServiceHandlers(store)
	router := handlers.NewRouter(serviceHandlers)

	return http.ListenAndServe(cfg.Address, middleware.RequestLogger(middleware.GzipMiddleware(router)))
}
