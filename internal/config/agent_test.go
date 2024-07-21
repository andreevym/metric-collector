package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAgentGetConfigFromFileSuccess(t *testing.T) {
	configPath := "../../config/agent.json"
	c := NewAgentConfig()
	err := c.GetConfigFromFile(configPath)
	require.NoError(t, err)
	require.Equal(t, c.Address, "localhost:8080")
	require.Equal(t, c.ReportInterval, 1)
	require.Equal(t, c.PollInterval, 1)
	require.Equal(t, c.LogLevel, "")
	require.Equal(t, c.SecretKey, "")
	require.Equal(t, c.RateLimit, 0)
	require.Equal(t, c.CryptoKey, "/path/to/key.pem")
}
