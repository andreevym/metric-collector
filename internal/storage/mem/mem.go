package mem

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"go.uber.org/zap"
)

type Storage struct {
	data map[string]*storage.Metric
	sync.RWMutex
	opt *BackupOptional
}

type BackupOptional struct {
	BackupPath    string
	StoreInterval time.Duration
}

func NewStorage(opt *BackupOptional) *Storage {
	return &Storage{
		map[string]*storage.Metric{},
		sync.RWMutex{},
		opt,
	}
}

func (s *Storage) Create(_ context.Context, m *storage.Metric) error {
	s.Lock()
	s.data[m.ID] = m
	s.Unlock()
	return nil
}

func (s *Storage) CreateAll(_ context.Context, metrics map[string]*storage.MetricR) error {
	s.Lock()
	for _, m := range metrics {
		s.data[m.Metric.ID] = m.Metric
	}
	s.Unlock()
	return nil
}

func (s *Storage) Read(_ context.Context, id string) (*storage.Metric, error) {
	v, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("%w: not found value by id %s", storage.ErrValueNotFound, id)
	}
	return v, nil
}

func (s *Storage) Update(_ context.Context, m *storage.Metric) error {
	s.Lock()
	_, ok := s.data[m.ID]
	if !ok {
		return fmt.Errorf(
			"can't update value by id, because value doesn't exists: id %s",
			m.ID,
		)
	}
	s.data[m.ID] = m
	s.Unlock()
	return nil
}

func (s *Storage) Delete(_ context.Context, id string) error {
	delete(s.data, id)
	return nil
}

func (s *Storage) Restore() error {
	if s.opt == nil || s.opt.BackupPath == "" {
		return nil
	}
	data, err := Load(s.opt.BackupPath)
	if err != nil {
		return err
	}
	s.data = data

	return nil
}

func (s *Storage) Backup() error {
	if s.opt == nil || s.opt.BackupPath == "" || s.opt.StoreInterval <= 0 {
		return nil
	}

	time.AfterFunc(s.opt.StoreInterval, func() {
		err := Save(s.opt.BackupPath, s.data)
		if err != nil {
			logger.Log.Error("problem to save backup ", zap.Error(err))
			panic(err)
		}
	})

	return nil
}
