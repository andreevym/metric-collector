package server

import (
	"log"
	"net/http"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/repository"
	"github.com/andreevym/metric-collector/internal/storage/mem"
)

type Server struct {
	counterMemStorage repository.Storage
	gaugeMemStorage   repository.Storage
}

func NewServer(counterMemStorage mem.Storage, gaugeMemStorage repository.Storage) Server {
	return Server{
		counterMemStorage: counterMemStorage,
		gaugeMemStorage:   gaugeMemStorage,
	}
}

func StartServer() {
	counterMemStorage := mem.NewStorage()
	gaugeMemStorage := mem.NewStorage()
	s := NewServer(counterMemStorage, gaugeMemStorage)

	mux := http.NewServeMux()
	mux.Handle("/update/", handlers.UpdateHandler(s.counterMemStorage, s.gaugeMemStorage))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
