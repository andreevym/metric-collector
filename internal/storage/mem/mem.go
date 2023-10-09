package mem

import (
	"fmt"
	"sync"
)

type Storage struct {
	m  map[string][]string
	rw sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		map[string][]string{},
		sync.RWMutex{},
	}
}

func (s *Storage) Create(key string, val string) error {
	s.rw.Lock()
	values := s.m[key]
	if values == nil {
		values = make([]string, 0)
	}
	s.m[key] = append(values, val)
	s.rw.Unlock()
	return nil
}

func (s *Storage) Read(key string) ([]string, error) {
	return s.m[key], nil
}

func (s *Storage) Update(key string, val []string) error {
	s.rw.Lock()
	if s.m[key] == nil {
		return fmt.Errorf("can't update value by key, because value doesn't exists: key %s",
			key)
	}
	s.m[key] = val
	s.rw.Unlock()
	return nil
}

func (s *Storage) Delete(key string) error {
	delete(s.m, key)
	return nil
}
