package handlers

import (
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/pg"
)

type ServiceHandlers struct {
	storage  *multistorage.Storage
	dbClient *pg.Client
}

func NewServiceHandlers(storage *multistorage.Storage, dbClient *pg.Client) *ServiceHandlers {
	return &ServiceHandlers{
		storage:  storage,
		dbClient: dbClient,
	}
}
