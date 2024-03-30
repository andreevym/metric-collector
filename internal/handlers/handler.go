package handlers

import (
	"github.com/andreevym/metric-collector/internal/storage/store"
)

type ServiceHandlers struct {
	storage  store.Storage
	dbClient store.Client
}

// NewServiceHandlers creates a new instance of ServiceHandlers with the provided dependencies.
func NewServiceHandlers(storage store.Storage, dbClient store.Client) *ServiceHandlers {
	return &ServiceHandlers{
		storage:  storage,
		dbClient: dbClient,
	}
}
