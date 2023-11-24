package pg

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Client struct {
	db *sql.DB
}

func NewClient(databaseDsn string) (*Client, error) {
	db, err := sql.Open("pgx", databaseDsn)
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

func (c *Client) Select(tableName string, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	q := fmt.Sprintf("SELECT value FROM %s WHERE key = $1;", tableName)
	rows := c.db.QueryRowContext(
		ctx,
		q,
		key,
	)
	if err := rows.Err(); err != nil {
		return "", err
	}
	var val string
	if err := rows.Scan(&val); err != nil {
		return "", err
	}
	return val, nil
}

func (c *Client) Insert(tableName string, key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	q := fmt.Sprintf("INSERT INTO %s (key, value) VALUES ($1, $2)", tableName)
	r, err := c.db.ExecContext(
		ctx,
		q,
		key,
		value,
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

func (c *Client) InsertAll(tableName string, kvMap map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed begin tx: %w", err)
	}

	stmt, err := tx.PrepareContext(
		ctx,
		fmt.Sprintf("INSERT INTO %s (key, value) VALUES ($1, $2)", tableName),
	)
	if err != nil {
		return fmt.Errorf("failed prepare context: %w", err)
	}
	defer stmt.Close()

	for k, v := range kvMap {
		_, err = stmt.ExecContext(
			ctx,
			k,
			v,
		)
		if err != nil {
			return fmt.Errorf("failed exec context: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed commit: %w", err)
	}

	return nil
}

func (c *Client) Update(
	tableName string,
	key string,
	value string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	q := fmt.Sprintf("UPDATE %s SET value = $1 WHERE key = $2", tableName)
	_, err := c.db.ExecContext(
		ctx,
		q,
		value,
		key,
	)
	if err != nil {
		return fmt.Errorf("failed update %w", err)
	}

	return nil
}

func (c *Client) Delete(tableName string, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	q := fmt.Sprintf("DELETE FROM %s WHERE key = $1", tableName)
	_, err := c.db.ExecContext(
		ctx,
		q,
		key,
	)
	if err != nil {
		return fmt.Errorf("failed delete: %w", err)
	}

	return nil
}

func (c *Client) ApplyMigration(sql string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := c.db.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("failed apply sql %s: %w", sql, err)
	}

	return nil
}
