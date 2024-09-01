package controller

import (
	"context"
	"fmt"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"go.uber.org/zap"
)

func (c Controller) Updates(ctx context.Context, metrics []*store.Metric) error {
	err := store.SaveAllMetric(ctx, c.storage, metrics)
	if err != nil {
		logger.Logger().Error("error updating metrics", zap.Error(err))
		return fmt.Errorf("error updating metrics: %w", err)
	}
	return nil
}
