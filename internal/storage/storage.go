package storage

import "context"

type Storage interface {
	CreateAll(ctx context.Context, metrics map[string]*MetricR) error
	Create(ctx context.Context, m *Metric) error
	Read(ctx context.Context, id string) (*Metric, error)
	Update(ctx context.Context, m *Metric) error
	Delete(ctx context.Context, id string) error
}

type MetricR struct {
	Metric   *Metric
	IsExists bool
}

type MType string

const (
	MTypeGauge   MType = "gauge"
	MTypeCounter MType = "counter"
)

type Metric struct {
	ID    string   `json:"id"`              // Имя метрики
	MType MType    `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}
