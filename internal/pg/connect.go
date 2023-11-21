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

func NewClient(databaseDsn string) *Client {
	db, err := sql.Open("pgx", databaseDsn)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}

	return &Client{db: db}
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
