package controller

import (
	"context"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"go.uber.org/zap"
)

func (c Controller) Value(ctx context.Context, id string, mType string) *store.Metric {
	foundMetric, err := c.storage.Read(ctx, id, mType)
	if err != nil {
		logger.Logger().Error(
			"failed to get value for metric",
			zap.String("id", id),
			zap.String("mType", mType),
			zap.Error(err),
		)
		return nil
	}
	return foundMetric
}
