package mem

import (
	"fmt"
	"sync"
)

type Storage struct {
	data map[string][]string
	rw   sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		map[string][]string{},
		sync.RWMutex{},
	}
}

func (s *Storage) Create(key string, val string) error {
	s.rw.Lock()
	values := s.data[key]
	if values == nil {
		values = make([]string, 0)
	}
	s.data[key] = append(values, val)
	s.rw.Unlock()
	return nil
}

func (s *Storage) Read(key string) ([]string, error) {
	v, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("%w: not found value by key %s", ErrValueNotFound, key)
	}
	return v, nil
}

func (s *Storage) Update(key string, val []string) error {
	s.rw.Lock()
	if s.data[key] == nil {
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

func (s *Storage) UpdateData(data map[string][]string) {
	s.data = data
}

func (s *Storage) Data() map[string][]string {
	return s.data
}
