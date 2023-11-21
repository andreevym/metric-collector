package server

import (
	"net/http"

	"github.com/andreevym/metric-collector/internal/config/serverconfig"
	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/pg"
	"github.com/andreevym/metric-collector/internal/storage/mem"
)

func Start(cfg *serverconfig.ServerConfig) error {
	counterMemStorage := mem.NewStorage()
	gaugeMemStorage := mem.NewStorage()
	store, err := multistorage.NewStorage(counterMemStorage, gaugeMemStorage, cfg)
	if err != nil {
		return err
	}

	dbClient, err := pg.NewClient(cfg.DatabaseDsn)
	if err != nil {
		return err
	}
	defer dbClient.Close()
	serviceHandlers := handlers.NewServiceHandlers(store, dbClient)
	router := handlers.NewRouter(
		serviceHandlers,
		middleware.GzipRequestMiddleware,
		middleware.GzipResponseMiddleware,
		middleware.RequestLogger,
	)

	return http.ListenAndServe(cfg.Address, router)
}
