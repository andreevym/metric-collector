package postgres

import (
	"sync"

	"github.com/andreevym/metric-collector/internal/pg"
)

type PgStorage struct {
	dbClient  *pg.Client
	tableName string
	rw        sync.RWMutex
}

func NewPgStorage(dbClient *pg.Client, tableName string) *PgStorage {
	return &PgStorage{
		dbClient,
		tableName,
		sync.RWMutex{},
	}
}

func (s *PgStorage) Create(key string, val string) error {
	s.rw.Lock()
	err := s.dbClient.Insert(s.tableName, key, val)
	s.rw.Unlock()
	return err
}

func (s *PgStorage) Read(key string) (string, error) {
	return s.dbClient.Select(s.tableName, key)
}

func (s *PgStorage) Update(key string, val string) error {
	s.rw.Lock()
	err := s.dbClient.Update(s.tableName, key, val)
	s.rw.Unlock()
	return err
}

func (s *PgStorage) Delete(key string) error {
	return s.dbClient.Delete(s.tableName, key)
}
