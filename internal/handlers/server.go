package handlers

import (
	"github.com/andreevym/metric-collector/internal/repository"
)

type Server struct {
	counterStorage repository.Storage
	gaugeStorage   repository.Storage
}

func NewServer(counterStorage repository.Storage, gaugeStorage repository.Storage) Server {
	return Server{
		counterStorage: counterStorage,
		gaugeStorage:   gaugeStorage,
	}
}

func (s Server) GaugeStorage() repository.Storage {
	return s.gaugeStorage
}

func (s Server) CounterStorage() repository.Storage {
	return s.counterStorage
}
