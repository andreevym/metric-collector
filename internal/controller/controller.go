package controller

import (
	"github.com/andreevym/metric-collector/internal/storage/store"
)

type Controller struct {
	storage  store.Storage
	dbClient store.Client
}

func NewController(storage store.Storage, dbClient store.Client) Controller {
	return Controller{
		storage:  storage,
		dbClient: dbClient,
	}
}
