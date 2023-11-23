package mem

import (
	"fmt"
	"sync"
	"time"

	"github.com/andreevym/metric-collector/internal/logger"
	"go.uber.org/zap"
)

type Storage struct {
	data map[string]string
	rw   sync.RWMutex
	opt  *BackupOptional
}

type BackupOptional struct {
	BackupPath    string
	StoreInterval time.Duration
}

func NewStorage(opt *BackupOptional) *Storage {
	return &Storage{
		map[string]string{},
		sync.RWMutex{},
		opt,
	}
}

func (s *Storage) Create(key string, val string) error {
	s.rw.Lock()
	s.data[key] = val
	s.rw.Unlock()
	return nil
}

func (s *Storage) Read(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("%w: not found value by key %s", ErrValueNotFound, key)
	}
	return v, nil
}

func (s *Storage) Update(key string, val string) error {
	s.rw.Lock()
	if s.data[key] == "" {
		return fmt.Errorf("can't update value by key, because value doesn't exists: key %s",
			key)
	}
	s.data[key] = val
	s.rw.Unlock()
	return nil
}

func (s *Storage) Delete(key string) error {
	delete(s.data, key)
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
