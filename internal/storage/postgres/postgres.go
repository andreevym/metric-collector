package postgres

import (
	"context"
	"sync"

	"github.com/andreevym/metric-collector/internal/storage"
)

type PgStorage struct {
	client *Client
	sync.RWMutex
}

func NewPgStorage(dbClient *Client) *PgStorage {
	return &PgStorage{
		dbClient,
		sync.RWMutex{},
	}
}

func (s *PgStorage) Create(ctx context.Context, m *storage.Metric) error {
	s.Lock()
	err := s.client.Insert(ctx, m)
	s.Unlock()
	return err
}

func (s *PgStorage) CreateAll(ctx context.Context, metrics map[string]*storage.MetricR) error {
	s.Lock()
	err := s.client.SaveAll(ctx, metrics)
	s.Unlock()
	return err
}

func (s *PgStorage) Read(ctx context.Context, id string, mType string) (*storage.Metric, error) {
	r, err := s.client.SelectByIDAndType(ctx, id, mType)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, storage.ErrValueNotFound
	}
	return r, err
}

func (s *PgStorage) Update(ctx context.Context, m *storage.Metric) error {
	s.Lock()
	err := s.client.Update(ctx, m)
	s.Unlock()
	return err
}

func (s *PgStorage) Delete(ctx context.Context, id string, mType string) error {
	return s.client.Delete(ctx, id, mType)
}
