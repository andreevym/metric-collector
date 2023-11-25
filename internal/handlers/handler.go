package handlers

import (
	"github.com/andreevym/metric-collector/internal/pg"
	"github.com/andreevym/metric-collector/internal/storage"
)

type ServiceHandlers struct {
	storage  storage.Storage
	dbClient *pg.Client
}

func NewServiceHandlers(storage storage.Storage, dbClient *pg.Client) *ServiceHandlers {
	return &ServiceHandlers{
		storage:  storage,
		dbClient: dbClient,
	}
}
