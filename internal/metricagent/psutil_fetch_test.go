package metricagent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	total, free, err := Memory()
	require.NoError(t, err)
	require.NotNil(t, total)
	require.NotNil(t, free)
}

func TestCPUutilization1(t *testing.T) {
	cpu, err := CPUUtilization()
	require.NoError(t, err)
	require.NotNil(t, cpu)
}
