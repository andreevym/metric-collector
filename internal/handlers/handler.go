package handlers

import (
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/pg"
)

type ServiceHandlers struct {
	metricStorage *multistorage.MetricStorage
	dbClient      *pg.Client
}

func NewServiceHandlers(storage *multistorage.MetricStorage, dbClient *pg.Client) *ServiceHandlers {
	return &ServiceHandlers{
		metricStorage: storage,
		dbClient:      dbClient,
	}
}
