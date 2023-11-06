package compressor_test

import (
	"testing"

	"github.com/andreevym/metric-collector/internal/compressor"
	"github.com/stretchr/testify/require"
)

func TestCompress(t *testing.T) {
	data := "test"
	compressedData, err := compressor.Compress([]byte(data))
	require.NoError(t, err)
	decompressedData, err := compressor.Decompress(compressedData)
	require.NoError(t, err)
	require.Equal(t, data, string(decompressedData))
}
