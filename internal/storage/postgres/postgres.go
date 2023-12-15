package postgres

import (
	"context"
	"errors"

	"github.com/andreevym/metric-collector/internal/logger"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/utils"
	"github.com/avast/retry-go"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

const (
	retryAttempts = 3
)

type PgStorage struct {
	client *Client
}

func NewPgStorage(dbClient *Client) *PgStorage {
	return &PgStorage{
		dbClient,
	}
}

func (s *PgStorage) Create(ctx context.Context, m *storage.Metric) error {
	var err error
	_ = retry.Do(
		func() error {
			err = s.client.Insert(ctx, m)
			if isRetriableError(err) {
				logger.Logger().Error("Retriable error detected. Retrying...", zap.Error(err))
				return err
			}
			return nil
		},
		retry.Attempts(retryAttempts),
		retry.DelayType(utils.RetryDelayType),
		retry.OnRetry(func(n uint, err error) {
			logger.Logger().Error("error send request to postgres",
				zap.Uint("currentAttempt", n),
				zap.Int("retryAttempts", retryAttempts),
				zap.Error(err),
			)
		}),
	)
	return err
}

func (s *PgStorage) CreateAll(ctx context.Context, metrics map[string]storage.MetricR) error {
	var err error
	_ = retry.Do(
		func() error {
			err = s.client.SaveAll(ctx, metrics)
			if isRetriableError(err) {
				logger.Logger().Error("Retriable error detected. Retrying...", zap.Error(err))
				return err
			}
			return nil
		},
		retry.Attempts(retryAttempts),
		retry.DelayType(utils.RetryDelayType),
		retry.OnRetry(func(n uint, err error) {
			logger.Logger().Error("error send request to postgres",
				zap.Uint("currentAttempt", n),
				zap.Int("retryAttempts", retryAttempts),
				zap.Error(err),
			)
		}),
	)
	return err
}

func (s *PgStorage) Read(ctx context.Context, id string, mType string) (*storage.Metric, error) {
	var m *storage.Metric
	var err error
	_ = retry.Do(
		func() error {
			m, err = s.client.SelectByIDAndType(ctx, id, mType)
			if isRetriableError(err) {
				logger.Logger().Error("Retriable error detected. Retrying...", zap.Error(err))
				return err
			}
			return nil
		},
		retry.Attempts(retryAttempts),
		retry.DelayType(utils.RetryDelayType),
		retry.OnRetry(func(n uint, err error) {
			logger.Logger().Error("error send request to postgres",
				zap.Uint("currentAttempt", n),
				zap.Int("retryAttempts", retryAttempts),
				zap.Error(err),
			)
		}),
	)
	return m, err
}

func (s *PgStorage) Update(ctx context.Context, m *storage.Metric) error {
	var err error
	_ = retry.Do(
		func() error {
			err = s.client.Update(ctx, m)
			if isRetriableError(err) {
				logger.Logger().Error("Retriable error detected. Retrying...", zap.Error(err))
				return err
			}
			return nil
		},
		retry.Attempts(retryAttempts),
		retry.DelayType(utils.RetryDelayType),
		retry.OnRetry(func(n uint, err error) {
			logger.Logger().Error("error send request to postgres",
				zap.Uint("currentAttempt", n),
				zap.Int("retryAttempts", retryAttempts),
				zap.Error(err),
			)
		}),
	)
	return err
}

func (s *PgStorage) Delete(ctx context.Context, id string, mType string) error {
	var err error
	_ = retry.Do(
		func() error {
			err = s.client.Delete(ctx, id, mType)
			if isRetriableError(err) {
				logger.Logger().Error("Retriable error detected. Retrying...", zap.Error(err))
				return err
			}
			return nil
		},
		retry.Attempts(retryAttempts),
		retry.DelayType(utils.RetryDelayType),
		retry.OnRetry(func(n uint, err error) {
			logger.Logger().Error("error send request to postgres",
				zap.Uint("currentAttempt", n),
				zap.Int("retryAttempts", retryAttempts),
				zap.Error(err),
			)
		}),
	)
	return err
}

func isRetriableError(err error) bool {
	var pgErr *pgconn.PgError
	// проверяем, что при обращении к PostgreSQL cервер получил ошибку транспорта
	// из категории Class 08 — Connection Exception.
	// если проблемы с соединением, то делаем повторяем попытку
	if err != nil && errors.As(err, &pgErr) &&
		pgerrcode.IsConnectionException(pgErr.Code) {
		return true
	}
	// если это ошибка не с соединением к PostgreSQL, то ретрай не нужен
	return false
}
