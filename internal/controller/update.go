package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"go.uber.org/zap"
)

func (c Controller) Update(ctx context.Context, metric *store.Metric) (*store.Metric, error) {
	if metric.MType != store.MTypeGauge && metric.MType != store.MTypeCounter {
		logger.Logger().Warn("unknown metric type", zap.String("type", metric.MType))
		return nil, fmt.Errorf("unknown metric type: %s", metric.MType)
	}

	foundValue, err := c.storage.Read(ctx, metric.ID, metric.MType)
	if err != nil && !errors.Is(err, store.ErrValueNotFound) {
		logger.Logger().Error("failed to read metric", zap.Error(err))
		return nil, fmt.Errorf("failed to read metric: %w", err)
	}

	if foundValue == nil {
		err = c.storage.Create(ctx, metric)
		if err != nil {
			logger.Logger().Error("failed create metric", zap.Error(err))
			return nil, fmt.Errorf("failed create metric: %w", err)
		}
	} else {
		if metric.MType == store.MTypeCounter {
			newDelta := *metric.Delta + *foundValue.Delta
			metric.Delta = &newDelta
		}
		err = c.storage.Update(ctx, metric)
		if err != nil {
			logger.Logger().Error("failed update metric", zap.Error(err))
			return nil, fmt.Errorf("failed update metric: %w", err)
		}
	}

	return metric, nil
}
