package mem

import (
	"os"
	"strconv"
	"testing"

	"github.com/andreevym/metric-collector/internal/storage/store"
	"github.com/stretchr/testify/require"
)

func TestEndToEndBackup(t *testing.T) {
	f, err := os.CreateTemp("", "tmpbackup")
	defer func() {
		err = os.RemoveAll(f.Name())
		require.NoError(t, err)
	}()

	data := make(map[string]*store.Metric)
	for i := 0; i < 1000; i++ {
		delta := int64(i)
		id := strconv.Itoa(i)
		data[id] = &store.Metric{
			ID:    id,
			MType: store.MTypeCounter,
			Delta: &delta,
			Value: nil,
		}
	}

	err = Save(f.Name(), data)
	require.NoError(t, err)

	loadedData, err := Load(f.Name())
	require.NoError(t, err)

	require.Equal(t, len(loadedData), len(data))

}
