package metricagent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIdentifyIP(t *testing.T) {
	ip, err := identifyIP()
	require.NoError(t, err)
	require.NotEmpty(t, ip)
}
