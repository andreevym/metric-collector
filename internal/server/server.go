package server

import (
	"net/http"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
)

func Start(address string) error {
	counterMemStorage := mem.NewStorage()
	gaugeMemStorage := mem.NewStorage()
	store, err := multistorage.NewStorage(counterMemStorage, gaugeMemStorage)
	if err != nil {
		return err
	}
	serviceHandlers := handlers.NewServiceHandlers(store)
	router := handlers.NewRouter(serviceHandlers)

	return http.ListenAndServe(address, middleware.RequestLogger(middleware.GzipMiddleware(router)))
}
