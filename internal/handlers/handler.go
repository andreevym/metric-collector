package handlers

import (
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/postgres"
)

type ServiceHandlers struct {
	storage  storage.Storage
	dbClient *postgres.Client
}

// NewServiceHandlers creates a new instance of ServiceHandlers with the provided dependencies.
func NewServiceHandlers(storage storage.Storage, dbClient *postgres.Client) *ServiceHandlers {
	return &ServiceHandlers{
		storage:  storage,
		dbClient: dbClient,
	}
}
