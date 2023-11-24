package gauge

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/require"
)

func TestGauge_Get_not_found(t *testing.T) {
	s := mem.NewStorage(nil)
	metricName := fmt.Sprintf("key_%f", rand.Float64())
	_, err := Get(s, metricName)
	require.ErrorIs(t, err, storage.ErrValueNotFound)
}

func TestGaugeEndToEnd(t *testing.T) {
	s := mem.NewStorage(nil)
	key1 := fmt.Sprintf("key_%f", rand.Float64())
	val11 := rand.Float64()
	store(t, key1, val11, s)
	get(t, key1, fmt.Sprintf("%v", val11), s)
	store(t, key1, val11, s)
	get(t, key1, fmt.Sprintf("%v", val11), s)
	val12 := rand.Float64()
	store(t, key1, val12, s)
	get(t, key1, fmt.Sprintf("%v", val12), s)
	val13 := rand.Float64()
	store(t, key1, val13, s)
	get(t, key1, fmt.Sprintf("%v", val13), s)
	key2 := fmt.Sprintf("key_%v", rand.Float64())
	val21 := rand.Float64()
	store(t, key2, val21, s)
	get(t, key1, fmt.Sprintf("%v", val13), s)
	get(t, key2, fmt.Sprintf("%v", val21), s)
	val22 := rand.Float64()
	store(t, key2, val22, s)
	get(t, key1, fmt.Sprintf("%v", val13), s)
	get(t, key2, fmt.Sprintf("%v", val22), s)
	val14 := rand.Float64()
	store(t, key1, val14, s)
	get(t, key1, fmt.Sprintf("%v", val14), s)
	get(t, key2, fmt.Sprintf("%v", val22), s)
}

func store(t *testing.T, metricName string, v float64, s *mem.Storage) {
	val1Str := fmt.Sprintf("%v", v)
	err := Validate(val1Str)
	require.NoError(t, err)
	_, err = Store(s, metricName, val1Str)
	require.NoError(t, err)
}

func get(t *testing.T, metricName string, expectedValue string, s *mem.Storage) {
	get, err := Get(s, metricName)
	require.NoError(t, err)
	require.Equal(t, expectedValue, get)
}
