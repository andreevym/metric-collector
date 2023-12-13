package storage

import (
	"context"
	"errors"
	"fmt"
)

type Storage interface {
	CreateAll(ctx context.Context, metrics map[string]MetricR) error
	Create(ctx context.Context, m *Metric) error
	Read(ctx context.Context, id string, mType string) (*Metric, error)
	Update(ctx context.Context, m *Metric) error
	Delete(ctx context.Context, id string, mType string) error
}

type MetricR struct {
	Metric   *Metric
	IsExists bool
}

const (
	MTypeGauge   string = "gauge"
	MTypeCounter string = "counter"
)

type Metric struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}

func SaveAllMetric(ctx context.Context, s Storage, metrics []*Metric) error {
	if len(metrics) == 0 {
		return nil
	}

	result := map[string]*Metric{}
	for _, metric := range metrics {
		found, ok := result[metric.ID+metric.MType]
		if ok && found != nil && metric.MType == MTypeCounter {
			newDelta := *metric.Delta + *found.Delta
			metric.Delta = &newDelta
		}

		result[metric.ID+metric.MType] = metric
	}

	metricsR := map[string]MetricR{}
	for _, metric := range metrics {
		found, err := s.Read(ctx, metric.ID, metric.MType)
		if err != nil && !errors.Is(err, ErrValueNotFound) {
			return fmt.Errorf("failed update metric: %w", err)
		}

		if found != nil && metric.MType == MTypeCounter {
			newDelta := *metric.Delta + *found.Delta
			metric.Delta = &newDelta
		}

		metricsR[metric.ID+metric.MType] = MetricR{
			Metric:   metric,
			IsExists: found != nil,
		}
	}

	err := s.CreateAll(ctx, metricsR)
	if err != nil {
		return fmt.Errorf("failed to create all metrics: %w", err)
	}
	return nil
}
