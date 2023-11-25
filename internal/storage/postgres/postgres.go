package postgres

import (
	"context"
	"sync"

	"github.com/andreevym/metric-collector/internal/pg"
	"github.com/andreevym/metric-collector/internal/storage"
)

type PgStorage struct {
	client *pg.Client
	sync.RWMutex
}

func NewPgStorage(dbClient *pg.Client) *PgStorage {
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
	err := s.client.InsertAll(ctx, metrics)
	s.Unlock()
	return err
}

func (s *PgStorage) Read(ctx context.Context, id string) (*storage.Metric, error) {
	r, err := s.client.SelectByID(ctx, id)
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

func (s *PgStorage) Delete(ctx context.Context, id string) error {
	return s.client.Delete(ctx, id)
}
