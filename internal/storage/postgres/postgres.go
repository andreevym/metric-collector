package postgres

import (
	"context"
	"errors"
	"sync"

	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/avast/retry-go"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
	var err error
	_ = retry.Do(func() error {
		err = s.client.Insert(ctx, m)
		if err != nil {
			var pgErr *pgconn.PgError
			// проверяем, что при обращении к PostgreSQL cервер получил ошибку транспорта
			// из категории Class 08 — Connection Exception.
			// если проблемы с соединением, то делаем повторяем попытку
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return err
			}
			// если это ошибка не с соединением к PostgreSQL, то ретрай не нужен
			return nil
		}
		return nil
	})
	s.Unlock()
	return err
}

func (s *PgStorage) CreateAll(ctx context.Context, metrics map[string]*storage.MetricR) error {
	s.Lock()
	var err error
	_ = retry.Do(func() error {
		err = s.client.SaveAll(ctx, metrics)
		if err != nil {
			var pgErr *pgconn.PgError
			// проверяем, что при обращении к PostgreSQL cервер получил ошибку транспорта
			// из категории Class 08 — Connection Exception.
			// если проблемы с соединением, то делаем повторяем попытку
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return err
			}
			// если это ошибка не с соединением к PostgreSQL, то ретрай не нужен
			return nil
		}
		return nil
	})
	s.Unlock()
	return err
}

func (s *PgStorage) Read(ctx context.Context, id string, mType string) (*storage.Metric, error) {
	var m *storage.Metric
	var err error
	_ = retry.Do(func() error {
		m, err = s.client.SelectByIDAndType(ctx, id, mType)
		if err != nil {
			var pgErr *pgconn.PgError
			if !errors.As(err, &pgErr) {
				// если это ошибка не с соединением к PostgreSQL, то ретрай не нужен
				return nil
			}
			if err.Error() == "sql: no rows in result set" {
				err = storage.ErrValueNotFound
				return nil
			}
			// проверяем, что при обращении к PostgreSQL cервер получил ошибку транспорта
			// из категории Class 08 — Connection Exception.
			// если проблемы с соединением, то делаем повторяем попытку
			if pgerrcode.IsConnectionException(pgErr.Code) {
				return err
			}
			// если это ошибка не с соединением к PostgreSQL, то ретрай не нужен
			return nil
		}
		return nil
	})
	return m, err
}

func (s *PgStorage) Update(ctx context.Context, m *storage.Metric) error {
	s.Lock()
	var err error
	_ = retry.Do(func() error {
		err = s.client.Update(ctx, m)
		if err != nil {
			var pgErr *pgconn.PgError
			// проверяем, что при обращении к PostgreSQL cервер получил ошибку транспорта
			// из категории Class 08 — Connection Exception.
			// если проблемы с соединением, то делаем повторяем попытку
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return err
			}
			// если это ошибка не с соединением к PostgreSQL, то ретрай не нужен
			return nil
		}
		return nil
	})
	s.Unlock()
	return err
}

func (s *PgStorage) Delete(ctx context.Context, id string, mType string) error {
	s.Lock()
	var err error
	_ = retry.Do(func() error {
		err = s.client.Delete(ctx, id, mType)
		if err != nil {
			var pgErr *pgconn.PgError
			// проверяем, что при обращении к PostgreSQL cервер получил ошибку транспорта
			// из категории Class 08 — Connection Exception.
			// если проблемы с соединением, то делаем повторяем попытку
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				return err
			}
			// если это ошибка не с соединением к PostgreSQL, то ретрай не нужен
			return nil
		}
		return nil
	})
	s.Unlock()
	return err
}
