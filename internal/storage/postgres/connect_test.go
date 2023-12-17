//go:build integration_test
// +build integration_test

package postgres_test

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/postgres"
	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	dsn := getDSN()
	pgClient, err := postgres.NewClient(dsn)
	require.NoError(t, err)
	defer pgClient.Close()

	err = pgClient.Ping()
	require.NoError(t, err)

	err = filepath.Walk("/home/yury/go/src/github.com/andreevym/metric-collector/migrations", func(path string, info fs.FileInfo, err error) error {
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

	id1 := "1"
	mType := storage.MTypeCounter

	delta := int64(1)

	insertedMetric1 := &storage.Metric{
		ID:    id1,
		MType: storage.MTypeCounter,
		Delta: &delta,
	}
	err = pgClient.Insert(context.TODO(), insertedMetric1)
	require.NoError(t, err)

	foundMetric, err := pgClient.SelectByIDAndType(context.TODO(), id1, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, insertedMetric1.Delta)

	delta = int64(2)
	updatedMetric1 := &storage.Metric{
		ID:    id1,
		MType: storage.MTypeCounter,
		Delta: &delta,
	}
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

	pgStorage := postgres.NewPgStorage(pgClient)

	err = pgStorage.Create(context.TODO(), insertedMetric1)
	require.NoError(t, err)

	foundMetric, err = pgStorage.Read(context.TODO(), id1, mType)
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

	id2 := "2"
	delta = int64(1)
	insertedMetric2 := &storage.Metric{
		ID:    id2,
		MType: storage.MTypeCounter,
		Delta: &delta,
	}
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

	delta = int64(2)
	updatedMetric2 := &storage.Metric{
		ID:    id2,
		MType: storage.MTypeCounter,
		Delta: &delta,
	}
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
}
