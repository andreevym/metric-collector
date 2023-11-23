package handlers

import (
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/pg"
)

type ServiceHandlers struct {
	metricStorage *multistorage.MetricManager
	dbClient      *pg.Client
}

func NewServiceHandlers(storage *multistorage.MetricManager, dbClient *pg.Client) *ServiceHandlers {
	return &ServiceHandlers{
		metricStorage: storage,
		dbClient:      dbClient,
	}
}
