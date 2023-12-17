//go:build integration_test
// +build integration_test

package postgres_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const (
	testDBName         = "test"
	testDBUserName     = "test"
	testDBUserPassword = "test"
)

var (
	getDSN          func() string
	getSUConnection func() (*pgx.Conn, error)
)

func initGetDSN(hostPort string) {
	getDSN = func() string {
		return fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			testDBUserName,
			testDBUserPassword,
			hostPort,
			testDBName,
		)
	}
}

func initSUConnection(ctx context.Context, hostPort string) error {
	getSUConnection = func() (*pgx.Conn, error) {
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
		return conn, err
	}
	return nil
}

func getHostPort(hostPort string) (string, uint16, error) {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("got an invalid host port string: %s", hostPort)
	}
	portStr := parts[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("failed to cast the port %s to an int: %w", portStr, err)
	}

	return parts[0], uint16(port), nil
}

func TestMain(m *testing.M) {
	code, err := runMain(m)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func runMain(m *testing.M) (int, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return 1, err
	}

	pg, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       "observer-integration-test",
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
		return 1, err
	}

	defer func() {
		if err := pool.Purge(pg); err != nil {
			log.Printf("failed to purge the postgres container: %v", err)
		}
	}()

	hostPort := pg.GetHostPort("5432/tcp")
	initGetDSN(hostPort)
	ctx := context.Background()
	if err := initSUConnection(ctx, hostPort); err != nil {
		return 1, err
	}

	pool.MaxWait = 10 * time.Second
	var conn *pgx.Conn
	if err := pool.Retry(func() error {
		conn, err = getSUConnection()
		if err != nil {
			return fmt.Errorf("failed to connect to the DB: %w", err)
		}
		return nil
	}); err != nil {
		return 1, err
	}

	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Printf("failed to correctly close the conenction: %v", err)
		}
	}()

	if err := createTestDB(ctx, conn); err != nil {
		return 1, fmt.Errorf("failed to create test db: %w", err)
	}

	exitCode := m.Run()
	return exitCode, nil
}

func createTestDB(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(
		ctx,
		fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'",
			testDBUserName,
			testDBUserPassword,
		),
	)

	if err != nil {
		return fmt.Errorf("failed to create a test user: %w", err)
	}

	_, err = conn.Exec(
		ctx,
		fmt.Sprintf(
			`CREATE DATABASE %s
			OWNER %s
			ENCODING 'UTF8'
			LC_COLLATE = 'en_US.utf8'
			LC_CTYPE = 'en_US.utf8'
			`,
			testDBName,
			testDBUserName,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create a test DB: %w", err)
	}

	return nil
}
