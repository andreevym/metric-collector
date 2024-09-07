package controller

import (
	"fmt"
	"github.com/andreevym/metric-collector/internal/logger"
	"go.uber.org/zap"
)

func (c Controller) Ping() error {
	if c.dbClient == nil {
		logger.Logger().Error("db client is nil")
		return fmt.Errorf("db client is nil")
	}
	err := c.dbClient.Ping()
	if err != nil {
		logger.Logger().Error("ping store", zap.Error(err))
		return fmt.Errorf("ping store: %w", err)
	}
	return nil
}
