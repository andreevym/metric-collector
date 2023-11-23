package mem

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEndToEndBackup(t *testing.T) {
	f, err := os.CreateTemp("", "tmpbackup")
	defer func() {
		err = os.RemoveAll(f.Name())
		require.NoError(t, err)
	}()

	data := make(map[string]string)
	for i := 0; i < 1000; i++ {
		formatInt := strconv.FormatInt(int64(i), 10)
		data[formatInt] = formatInt
	}

	err = Save(f.Name(), data)
	require.NoError(t, err)

	loadedData, err := Load(f.Name())
	require.NoError(t, err)

	require.Equal(t, len(loadedData), len(data))

}
