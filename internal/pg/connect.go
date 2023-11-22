package pg

import (
	"context"
	"database/sql"
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

	row := c.db.QueryRowContext(ctx, "SELECT value FROM ? WHERE key = ?", tableName, key)
	var res string
	err := row.Scan(&res)
	if err != nil {
		return "", err
	}
	err = row.Err()
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *Client) Insert(tableName string, key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := c.db.ExecContext(ctx, "INSERT INTO ? (key, value) VALUES (?, ?)", tableName, key, value)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Update(tableName string, key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := c.db.ExecContext(ctx, "INSERT INTO ? (key, value) VALUES (? ?)", tableName, key, value)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Delete(tableName string, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := c.db.ExecContext(ctx, "DELETE FROM ? WHERE key = ?", tableName, key)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ApplyMigration(sql string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := c.db.ExecContext(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}
