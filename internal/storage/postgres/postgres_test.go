//go:build integration_test

package postgres_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	testDBUserName     = "test"
	testDBUserPassword = "test"
	containerName      = "integration-test-postgres"
)

var (
	pool     *dockertest.Pool
	pg       *dockertest.Resource
	hostPort string
)

func getDSN(
	hostPort string,
	dbName string,
	dbUserName string,
	dbUserPassword string,
) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		dbUserName,
		dbUserPassword,
		hostPort,
		dbName,
	)
}

func getSUConnection(ctx context.Context, hostPort string) (*pgx.Conn, error) {
	dsnPostgres := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		"postgres",
		"postgres",
		hostPort,
		"postgres",
	)
	conn, err := pgx.Connect(ctx, dsnPostgres)
	if err != nil {
		return nil, fmt.Errorf("failed to get a super user connection: %w", err)
	}
	return conn, nil
}

func downPostgres() {
	pool, err := dockertest.NewPool("")
	if err != nil {
		panic(fmt.Errorf("failed to create new pool for dockertest: %w", err))
	}

	err = pool.RemoveContainerByName(containerName)
	if err != nil {
		panic(fmt.Errorf("failed to remove container by name '%s': %w", containerName, err))
	}
}

func TestMain(m *testing.M) {
	code, err := runMain(m)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func runMain(m *testing.M) (int, error) {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		return 1, fmt.Errorf("failed to create new pool for dockertest: %w", err)
	}

	err = pool.RemoveContainerByName(containerName)
	if err != nil {
		return 1, fmt.Errorf("failed to remove container by name '%s': %w", containerName, err)
	}

	pg, err = pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       containerName,
			Repository: "postgres",
			Tag:        "15.3",
			Env: []string{
				"POSTGRES_USER=postgres",
				"POSTGRES_PASSWORD=postgres",
			},
			ExposedPorts: []string{"5432"},
		},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return 1, fmt.Errorf("failed to run postgres container with integration test: %w", err)
	}
	hostPort = pg.GetHostPort("5432/tcp")
	pool.MaxWait = 10 * time.Second

	err = createTestUser(context.Background(), testDBUserName, testDBUserPassword)
	if err != nil {
		return 1, fmt.Errorf("failed to create test user: %w", err)
	}

	exitCode := m.Run()
	downPostgres()
	return exitCode, nil
}

func createTestUser(
	ctx context.Context,
	dbUserName string,
	dbUserPassword string,
) error {
	var err error
	var conn *pgx.Conn
	if err = pool.Retry(func() error {
		conn, err = getSUConnection(ctx, hostPort)
		if err != nil {
			return fmt.Errorf("failed to connect to the DB: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to get su connection: %w", err)
	}

	defer func() {
		if err = conn.Close(ctx); err != nil {
			log.Printf("failed to correctly close the conenction: %v", err)
		}
	}()

	_, err = conn.Exec(
		ctx,
		fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'",
			dbUserName,
			dbUserPassword,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create a test user: %w", err)
	}
	return nil
}

func CreateTestDB(
	ctx context.Context,
	dbName string,
	dbUserName string,
) error {
	var err error
	var conn *pgx.Conn
	if err = pool.Retry(func() error {
		conn, err = getSUConnection(ctx, hostPort)
		if err != nil {
			return fmt.Errorf("failed to connect to the DB: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to get su connection: %w", err)
	}

	defer func() {
		if err = conn.Close(ctx); err != nil {
			log.Printf("failed to correctly close the conenction: %v", err)
		}
	}()

	_, err = conn.Exec(
		ctx,
		fmt.Sprintf(
			`CREATE DATABASE %s
			OWNER %s
			ENCODING 'UTF8'
			LC_COLLATE = 'en_US.utf8'
			LC_CTYPE = 'en_US.utf8'
			`,
			dbName,
			dbUserName,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create a test db: %w", err)
	}

	return nil
}

func DropTestDB(
	ctx context.Context,
	dbName string,
) error {
	var err error
	var conn *pgx.Conn
	if err = pool.Retry(func() error {
		conn, err = getSUConnection(ctx, hostPort)
		if err != nil {
			return fmt.Errorf("failed to connect to the DB: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to get su connection: %w", err)
	}

	defer func() {
		if err = conn.Close(ctx); err != nil {
			log.Printf("failed to correctly close the conenction: %v", err)
		}
	}()

	_, err = conn.Exec(
		ctx,
		fmt.Sprintf("DROP DATABASE %s", dbName),
	)
	if err != nil {
		return fmt.Errorf("failed to drop a test db: %w", err)
	}

	return nil
}
