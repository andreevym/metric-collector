package gauge

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/require"
)

func TestGauge_Get_not_found(t *testing.T) {
	storage := mem.NewStorage()
	metricName := fmt.Sprintf("key_%f", rand.Float64())
	_, err := Get(storage, metricName)
	require.ErrorIs(t, err, mem.ErrValueNotFound)
}

func TestGaugeEndToEnd(t *testing.T) {
	storage := mem.NewStorage()
	key1 := fmt.Sprintf("key_%f", rand.Float64())
	val11 := rand.Float64()
	store(t, key1, val11, storage)
	get(t, key1, []string{
		fmt.Sprintf("%v", val11),
	}, storage)
	store(t, key1, val11, storage)
	get(t, key1, []string{
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val11),
	}, storage)
	val12 := rand.Float64()
	store(t, key1, val12, storage)
	get(t, key1, []string{
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val12),
	}, storage)
	val13 := rand.Float64()
	store(t, key1, val13, storage)
	get(t, key1, []string{
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val12),
		fmt.Sprintf("%v", val13),
	}, storage)
	key2 := fmt.Sprintf("key_%v", rand.Float64())
	val21 := rand.Float64()
	store(t, key2, val21, storage)
	get(t, key1, []string{
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val12),
		fmt.Sprintf("%v", val13),
	}, storage)
	get(t, key2, []string{
		fmt.Sprintf("%v", val21),
	}, storage)
	val22 := rand.Float64()
	store(t, key2, val22, storage)
	get(t, key1, []string{
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val12),
		fmt.Sprintf("%v", val13),
	}, storage)
	get(t, key2, []string{
		fmt.Sprintf("%v", val21),
		fmt.Sprintf("%v", val22),
	}, storage)
	val14 := rand.Float64()
	store(t, key1, val14, storage)
	get(t, key1, []string{
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val11),
		fmt.Sprintf("%v", val12),
		fmt.Sprintf("%v", val13),
		fmt.Sprintf("%v", val14),
	}, storage)
	get(t, key2, []string{
		fmt.Sprintf("%v", val21),
		fmt.Sprintf("%v", val22),
	}, storage)
}

func store(t *testing.T, metricName string, v float64, storage *mem.Storage) {
	val1Str := fmt.Sprintf("%v", v)
	err := Validate(val1Str)
	require.NoError(t, err)
	err = Store(storage, metricName, val1Str)
	require.NoError(t, err)
}

func get(t *testing.T, metricName string, expectedValue []string, storage *mem.Storage) {
	get, err := Get(storage, metricName)
	require.NoError(t, err)
	require.Equal(t, expectedValue, get)
}
