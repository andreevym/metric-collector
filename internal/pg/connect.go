package pg

import (
	"context"
	"database/sql"
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

func (c *Client) SelectByID(ctx context.Context, id string) (*storage.Metric, error) {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	metric := storage.Metric{}
	err := c.db.GetContext(rCtx, &metric, "SELECT * FROM metric WHERE id = $1;", id)
	if err != nil {
		return nil, fmt.Errorf("failed execute select: %w", err)
	}

	return &metric, nil
}

func (c *Client) Insert(ctx context.Context, metrics *storage.Metric) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	r, err := c.db.ExecContext(
		rCtx,
		"INSERT INTO metric (id, type, delta, value) VALUES (@id, @type, @delta, @value)",
		sql.Named("id", metrics.ID),
		sql.Named("type", metrics.MType),
		sql.Named("delta", metrics.Delta),
		sql.Named("value", metrics.Value),
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

func (c *Client) InsertAll(ctx context.Context, metrics map[string]*storage.MetricR) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}

	insStmt, err := tx.PrepareContext(
		rCtx,
		"INSERT INTO metric (id, type, delta, value) VALUES (@id, @type, @delta, @value)",
	)
	if err != nil {
		return fmt.Errorf("failed prepare context: %w", err)
	}
	defer insStmt.Close()

	updStmt, err := tx.PrepareContext(
		rCtx,
		"UPDATE metric SET delta = @delta, value = @value WHERE id = @id",
	)
	if err != nil {
		return fmt.Errorf("failed prepare context: %w", err)
	}
	defer updStmt.Close()

	for _, m := range metrics {
		if m.IsExists {
			_, err = updStmt.ExecContext(
				rCtx,
				sql.Named("id", m.Metric.ID),
				sql.Named("delta", m.Metric.Delta),
				sql.Named("value", m.Metric.Value),
			)
			if err != nil {
				return fmt.Errorf("failed exec context: %w", err)
			}
		} else {
			_, err = insStmt.ExecContext(
				rCtx,
				sql.Named("id", m.Metric.ID),
				sql.Named("type", m.Metric.MType),
				sql.Named("delta", m.Metric.Delta),
				sql.Named("value", m.Metric.Value),
			)
			if err != nil {
				return fmt.Errorf("failed exec context: %w", err)
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
	metric *storage.Metric,
) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := c.db.ExecContext(
		rCtx,
		"UPDATE metric SET delta = @delta, value = @value WHERE id = @id",
		sql.Named("id", metric.ID),
		sql.Named("delta", metric.Delta),
		sql.Named("value", metric.Value),
	)
	if err != nil {
		return fmt.Errorf("failed update %w", err)
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	rCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	_, err := c.db.ExecContext(rCtx, "DELETE FROM metric WHERE key = $1", id)
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
