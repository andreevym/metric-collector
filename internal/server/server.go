package server

import (
	"log"
	"net/http"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/storage/mem"
)

func StartServer(addr string) {
	counterMemStorage := mem.NewStorage()
	gaugeMemStorage := mem.NewStorage()
	s := handlers.NewServer(counterMemStorage, gaugeMemStorage)
	log.Fatal(http.ListenAndServe(addr, handlers.Router(s)))
}
