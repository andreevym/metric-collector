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

	id := "123"
	mType := storage.MTypeCounter

	delta := int64(1)

	insertedMetric := &storage.Metric{
		ID:    id,
		MType: storage.MTypeCounter,
		Delta: &delta,
	}
	err = pgClient.Insert(context.TODO(), insertedMetric)
	require.NoError(t, err)

	foundMetric, err := pgClient.SelectByIDAndType(context.TODO(), id, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, insertedMetric.Delta)

	delta = int64(2)
	updatedMetric := &storage.Metric{
		ID:    id,
		MType: storage.MTypeCounter,
		Delta: &delta,
	}
	err = pgClient.Update(context.TODO(), updatedMetric)
	require.NoError(t, err)

	foundMetric, err = pgClient.SelectByIDAndType(context.TODO(), id, mType)
	require.NoError(t, err)
	require.NotNil(t, foundMetric)
	require.Equal(t, foundMetric.Delta, updatedMetric.Delta)

	err = pgClient.Delete(context.TODO(), id, mType)
	require.NoError(t, err)

	foundMetric, err = pgClient.SelectByIDAndType(context.TODO(), id, mType)
	require.EqualError(t, err, "not found value")
	require.Nil(t, foundMetric)
}
