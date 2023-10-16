package handlers

import "github.com/andreevym/metric-collector/internal/multistorage"

type ServiceHandlers struct {
	storage *multistorage.Storage
}

func NewServiceHandlers(storage *multistorage.Storage) *ServiceHandlers {
	return &ServiceHandlers{
		storage: storage,
	}
}
