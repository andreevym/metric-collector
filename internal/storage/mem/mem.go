package mem

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"go.uber.org/zap"
)

type Storage struct {
	data map[string]*store.Metric
	sync.RWMutex
	opt *BackupOptional
}

type BackupOptional struct {
	BackupPath    string
	StoreInterval time.Duration
}

func NewStorage(opt *BackupOptional) *Storage {
	return &Storage{
		map[string]*store.Metric{},
		sync.RWMutex{},
		opt,
	}
}

func (s *Storage) Create(_ context.Context, m *store.Metric) error {
	s.Lock()
	if m.MType != store.MTypeGauge && m.MType != store.MTypeCounter {
		return fmt.Errorf("metric type %s is not valid for ID %s", m.MType, m.ID)
	}
	s.data[m.ID+m.MType] = m
	s.Unlock()
	err := s.Backup()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateAll(_ context.Context, metrics map[string]store.MetricR) error {
	s.Lock()
	for _, m := range metrics {
		if m.Metric.MType != store.MTypeGauge && m.Metric.MType != store.MTypeCounter {
			return fmt.Errorf("metric type %s is not valid for ID %s", m.Metric.MType, m.Metric.ID)
		}
		s.data[m.Metric.ID+m.Metric.MType] = m.Metric
	}
	s.Unlock()
	err := s.Backup()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Read(_ context.Context, id string, mType string) (*store.Metric, error) {
	v, ok := s.data[id+mType]
	if !ok {
		return nil, fmt.Errorf("%w: not found value by id %s", store.ErrValueNotFound, id)
	}
	return v, nil
}

func (s *Storage) Update(_ context.Context, m *store.Metric) error {
	s.Lock()
	if m.MType != store.MTypeGauge && m.MType != store.MTypeCounter {
		return fmt.Errorf("metric type %s is not valid for ID %s", m.MType, m.ID)
	}
	_, ok := s.data[m.ID+m.MType]
	if !ok {
		return fmt.Errorf(
			"can't update value by id, because value doesn't exists: id %s",
			m.ID,
		)
	}
	s.data[m.ID+m.MType] = m
	s.Unlock()
	err := s.Backup()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Delete(_ context.Context, id string, mType string) error {
	delete(s.data, id+mType)
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
	if s.opt == nil || s.opt.BackupPath == "" || s.opt.StoreInterval < 0 {
		return nil
	}

	err := Save(s.opt.BackupPath, s.data)
	if err != nil {
		logger.Logger().Error("problem to save backup ", zap.Error(err))
		return fmt.Errorf("save backup: %s", err)
	}

	return nil
}
