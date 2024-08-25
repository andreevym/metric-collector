package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerGetConfigFromFileSuccess(t *testing.T) {
	configPath := "../../config/server.json"
	c := NewServerConfig()
	err := c.GetConfigFromFile(configPath)
	require.NoError(t, err)
	require.Equal(t, c.Address, "localhost:8080")
	require.Equal(t, c.LogLevel, "")
	require.Equal(t, c.StoreInterval, 1)
	require.Equal(t, c.FileStoragePath, "/path/to/file.db")
	require.Equal(t, c.Restore, true)
	require.Equal(t, c.DatabaseDsn, "")
	require.Equal(t, c.SecretKey, "")
	require.Equal(t, c.CryptoKey, "/path/to/key.pem")
}
