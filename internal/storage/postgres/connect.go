package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andreevym/metric-collector/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	db *sqlx.DB
}

func NewClient(databaseDsn string) (*Client, error) {
	db, err := sqlx.Open("pgx", databaseDsn)
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := c.db.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SelectByIDAndType(ctx context.Context, id string, mType string) (*storage.Metric, error) {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	metrics := []storage.Metric{}
	err := c.db.SelectContext(
		rCtx,
		&metrics,
		"SELECT id as \"id\", type as \"mtype\", delta as \"delta\", value as \"value\" "+
			"FROM metric WHERE id = $1 and type = $2;",
		id,
		mType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed execute select: %w", err)
	}
	if len(metrics) == 0 {
		return nil, storage.ErrValueNotFound
	}
	if len(metrics) > 1 {
		return nil, fmt.Errorf("something goung wrong, expect single value, but found %d", len(metrics))
	}

	return &metrics[0], nil
}

func (c *Client) Insert(ctx context.Context, m *storage.Metric) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if m.MType != storage.MTypeGauge && m.MType != storage.MTypeCounter {
		return fmt.Errorf("metric type %s is not valid for ID %s", m.MType, m.ID)
	}

	if m.Delta == nil && m.Value == nil {
		return errors.New("metric can't have nil delta and nil value")
	}
	r, err := c.db.ExecContext(
		rCtx,
		"INSERT INTO metric (id, type, delta, value) VALUES ($1, $2, $3, $4)",
		m.ID,
		m.MType,
		m.Delta,
		m.Value,
	)
	if err != nil {
		return fmt.Errorf("failed insert %w", err)
	}
	_, err = r.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed insert %w", err)
	}

	return nil
}

func (c *Client) SaveAll(ctx context.Context, metrics map[string]storage.MetricR) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}

	insStmt, err := tx.PrepareContext(
		rCtx,
		"INSERT INTO metric (id, type, delta, value) VALUES ($1, $2, $3, $4)",
	)
	if err != nil {
		return fmt.Errorf("failed prepare context: %w", err)
	}
	defer insStmt.Close()

	updStmt, err := tx.PrepareContext(
		rCtx,
		"UPDATE metric SET delta = $2, value = $3 WHERE id = $1 and type = $4",
	)
	if err != nil {
		return fmt.Errorf("failed prepare context: %w", err)
	}
	defer updStmt.Close()

	for _, m := range metrics {
		if m.Metric.MType != storage.MTypeGauge && m.Metric.MType != storage.MTypeCounter {
			return fmt.Errorf("metric type %s is not valid for ID %s", m.Metric.MType, m.Metric.ID)
		}

		if m.Metric.Delta == nil && m.Metric.Value == nil {
			return errors.New("metric can't have nil delta and nil value")
		}

		if m.IsExists {
			_, err = updStmt.ExecContext(
				rCtx,
				m.Metric.ID,
				m.Metric.Delta,
				m.Metric.Value,
				m.Metric.MType,
			)
			if err != nil {
				return fmt.Errorf("failed update: %w", err)
			}
		} else {
			_, err = insStmt.ExecContext(
				rCtx,
				m.Metric.ID,
				m.Metric.MType,
				m.Metric.Delta,
				m.Metric.Value,
			)
			if err != nil {
				return fmt.Errorf("failed insert: %w", err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed commit: %w", err)
	}

	return nil
}

func (c *Client) Update(
	ctx context.Context,
	m *storage.Metric,
) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if m.MType != storage.MTypeGauge && m.MType != storage.MTypeCounter {
		return fmt.Errorf("metric type %s is not valid for ID %s", m.MType, m.ID)
	}

	if m.Delta == nil && m.Value == nil {
		return errors.New("metric can't have nil delta and nil value")
	}
	_, err := c.db.ExecContext(
		rCtx,
		"UPDATE metric SET delta = $2, value = $3 WHERE id = $1 and type = $4",
		m.ID,
		m.Delta,
		m.Value,
		m.MType,
	)
	if err != nil {
		return fmt.Errorf("failed update %w", err)
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, id string, mType string) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := c.db.ExecContext(rCtx, "DELETE FROM metric WHERE key = $1 and type = $2", id, mType)
	if err != nil {
		return fmt.Errorf("failed delete: %w", err)
	}

	return nil
}

func (c *Client) ApplyMigration(ctx context.Context, sql string) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := c.db.ExecContext(rCtx, sql)
	if err != nil {
		return fmt.Errorf("failed apply sql '%s': %w", sql, err)
	}

	return nil
}
