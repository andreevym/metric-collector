// Package storage provides interfaces and implementations for metric store.
package store

import (
	"context"
	"errors"
	"fmt"
)

// ErrValueNotFound indicates that the requested metric value was not found.
var ErrValueNotFound = errors.New("not found value")

// Storage defines the interface for metric storage operations.
//
//go:generate mockgen -destination=../mocks/mock_store.go -source=store.go -package=mocks Storage
type Storage interface {
	CreateAll(ctx context.Context, metrics map[string]MetricR) error
	Create(ctx context.Context, m *Metric) error
	Read(ctx context.Context, id string, mType string) (*Metric, error)
	Update(ctx context.Context, m *Metric) error
	Delete(ctx context.Context, id string, mType string) error
	Backup() error
}

//go:generate mockgen -destination=../postgres/mock_pgclient.go -source=store.go -package=postgres Client
type Client interface {
	Close() error
	Ping() error
	SelectByIDAndType(ctx context.Context, id string, mType string) (*Metric, error)
	Insert(ctx context.Context, m *Metric) error
	SaveAll(ctx context.Context, metrics map[string]MetricR) error
	Update(context.Context, *Metric) error
	Delete(context.Context, string, string) error
	ApplyMigration(ctx context.Context, sql string) error
}

// MetricR represents a metric along with a flag indicating its existence in the store.
type MetricR struct {
	Metric   *Metric // Metric information
	IsExists bool    // Flag indicating whether the metric already exists in the storage
}

// Metric represents a metric with its ID, type, delta, and value.
type Metric struct {
	ID    string   `json:"id"`              // Metric ID
	MType string   `json:"type"`            // Metric type: gauge or counter
	Delta *int64   `json:"delta,omitempty"` // Delta value (applicable for counter type)
	Value *float64 `json:"value,omitempty"` // Value (applicable for gauge type)
}

// MType constants represent different metric types.
const (
	MTypeGauge   string = "gauge"
	MTypeCounter string = "counter"
)

// SaveAllMetric saves multiple metrics in the store.
// It takes a context, a storage instance, and a slice of Metric pointers.
// Metrics are grouped by their ID and type to avoid duplication.
// If a metric with the same ID and type already exists in the storage and its type is counter,
// the delta of the existing metric and the new metric are summed up.
// After processing the metrics, they are saved in the storage using the CreateAll method.
// If any error occurs during the process, an error is returned.
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
