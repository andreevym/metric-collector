package server

import (
	"log"
	"net/http"

	"github.com/andreevym/metric-collector/internal/storage/mem"
)

type Storage interface {
	Create(key string, val string) error
	Read(key string) ([]string, error)
	Update(key string, val []string) error
	Delete(key string) error
}

type Server struct {
	counterMemStorage Storage
	gaugeMemStorage   Storage
}

func NewServer(counterMemStorage mem.Storage, gaugeMemStorage Storage) Server {
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
	mux.Handle("/update/", http.HandlerFunc(s.update))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
