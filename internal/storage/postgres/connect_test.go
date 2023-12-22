//go:build integration_test

package postgres_test

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/postgres"
	"github.com/stretchr/testify/require"
)

// test data
var (
	id1             = "1"
	delta1          = int64(1)
	delta2          = int64(2)
	mType           = storage.MTypeCounter
	insertedMetric1 = &storage.Metric{
		ID:    id1,
		MType: mType,
		Delta: &delta1,
	}
	updatedMetric1 = &storage.Metric{
		ID:    id1,
		MType: mType,
		Delta: &delta2,
	}
	id2             = "2"
	insertedMetric2 = &storage.Metric{
		ID:    id2,
		MType: storage.MTypeCounter,
		Delta: &delta1,
	}
	updatedMetric2 = &storage.Metric{
		ID:    id2,
		MType: storage.MTypeCounter,
		Delta: &delta2,
	}
)

func TestPgStorageEndToEnd(t *testing.T) {
	defer func() {
		if rc := recover(); rc != nil {
			msg := fmt.Sprintf("panic: %v", rc)
			require.Fail(t, msg)
		}
	}()

	ctx := context.Background()
	dbName := strings.ToLower(t.Name())
	err := CreateTestDB(ctx, dbName, testDBUserName)
	require.NoError(t, err)

	dsn := getDSN(hostPort, dbName, testDBUserName, testDBUserPassword)
	pgClient, err := postgres.NewClient(dsn)
	require.NoError(t, err)

	err = pgClient.Ping()
	require.NoError(t, err)

	migrate(t, pgClient)

	pgStorage := postgres.NewPgStorage(pgClient)

	err = pgStorage.Create(context.TODO(), insertedMetric1)
	require.NoError(t, err)

	foundMetric, err := pgStorage.Read(context.TODO(), id1, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, insertedMetric1.Delta)

	err = pgStorage.Update(context.TODO(), updatedMetric1)
	require.NoError(t, err)

	foundMetric, err = pgStorage.Read(context.TODO(), id1, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, updatedMetric1.Delta)

	err = pgStorage.Delete(context.TODO(), id1, mType)
	require.NoError(t, err)

	foundMetric, err = pgStorage.Read(context.TODO(), id1, mType)
	require.EqualError(t, err, "not found value")
	require.Nil(t, foundMetric)

	createdMetrics := map[string]storage.MetricR{}
	createdMetrics[id1] = storage.MetricR{
		Metric:   insertedMetric1,
		IsExists: false,
	}
	createdMetrics[id2] = storage.MetricR{
		Metric:   insertedMetric2,
		IsExists: false,
	}
	err = pgStorage.CreateAll(context.TODO(), createdMetrics)
	require.NoError(t, err)

	foundMetric, err = pgStorage.Read(context.TODO(), id1, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, insertedMetric1.Delta)

	foundMetric, err = pgStorage.Read(context.TODO(), id2, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, insertedMetric2.Delta)

	updatedMetrics := map[string]storage.MetricR{}
	updatedMetrics[id1] = storage.MetricR{
		Metric:   updatedMetric1,
		IsExists: true,
	}
	updatedMetrics[id2] = storage.MetricR{
		Metric:   updatedMetric2,
		IsExists: true,
	}
	err = pgStorage.CreateAll(context.TODO(), updatedMetrics)
	require.NoError(t, err)

	foundMetric, err = pgStorage.Read(context.TODO(), id1, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, updatedMetric1.Delta)

	foundMetric, err = pgStorage.Read(context.TODO(), id2, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, updatedMetric2.Delta)

	err = pgClient.Close()
	require.NoError(t, err)
	err = DropTestDB(ctx, dbName)
	require.NoError(t, err)
}

func TestPgClientEndToEnd(t *testing.T) {
	defer func() {
		if rc := recover(); rc != nil {
			msg := fmt.Sprintf("panic: %v", rc)
			require.Fail(t, msg)
		}
	}()

	ctx := context.Background()
	dbName := strings.ToLower(t.Name())
	err := CreateTestDB(ctx, dbName, testDBUserName)
	require.NoError(t, err)

	dsn := getDSN(hostPort, dbName, testDBUserName, testDBUserPassword)
	pgClient, err := postgres.NewClient(dsn)
	require.NoError(t, err)

	err = pgClient.Ping()
	require.NoError(t, err)

	migrate(t, pgClient)

	err = pgClient.Insert(context.TODO(), insertedMetric1)
	require.NoError(t, err)

	foundMetric, err := pgClient.SelectByIDAndType(context.TODO(), id1, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, insertedMetric1.Delta)

	err = pgClient.Update(context.TODO(), updatedMetric1)
	require.NoError(t, err)

	foundMetric, err = pgClient.SelectByIDAndType(context.TODO(), id1, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, updatedMetric1.Delta)

	err = pgClient.Delete(context.TODO(), id1, mType)
	require.NoError(t, err)

	foundMetric, err = pgClient.SelectByIDAndType(context.TODO(), id1, mType)
	require.EqualError(t, err, "not found value")
	require.Nil(t, foundMetric)

	err = pgClient.Close()
	require.NoError(t, err)
	err = DropTestDB(ctx, dbName)
	require.NoError(t, err)
}

func migrate(t *testing.T, pgClient *postgres.Client) {
	err := filepath.Walk("../../../migrations", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			bytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			err = pgClient.ApplyMigration(context.TODO(), string(bytes))
			if err != nil {
				return err
			}
		}

		return nil
	})
	require.NoError(t, err)
}
