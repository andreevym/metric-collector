package server

import (
	"log"
	"net/http"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
)

func Start(address string) {
	counterMemStorage := mem.NewStorage()
	gaugeMemStorage := mem.NewStorage()
	store, err := multistorage.NewStorage(counterMemStorage, gaugeMemStorage)
	if err != nil {
		log.Fatal(err)
	}
	serviceHandlers := handlers.NewServiceHandlers(store)
	router := handlers.NewRouter(serviceHandlers)
	log.Fatal(http.ListenAndServe(address, router))
}
