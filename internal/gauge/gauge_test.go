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
	metricName := fmt.Sprintf("key_%d", rand.Int())
	_, err := Get(storage, metricName)
	require.ErrorIs(t, err, mem.ErrValueNotFound)
}

func TestGaugeEndToEnd(t *testing.T) {
	storage := mem.NewStorage()
	key1 := fmt.Sprintf("key_%d", rand.Int())
	val11 := rand.Int()
	store(t, key1, val11, storage)
	get(t, key1, []string{
		fmt.Sprintf("%d", val11),
	}, storage)
	store(t, key1, val11, storage)
	get(t, key1, []string{
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val11),
	}, storage)
	val12 := rand.Int()
	store(t, key1, val12, storage)
	get(t, key1, []string{
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val12),
	}, storage)
	val13 := rand.Int()
	store(t, key1, val13, storage)
	get(t, key1, []string{
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val12),
		fmt.Sprintf("%d", val13),
	}, storage)
	key2 := fmt.Sprintf("key_%d", rand.Int())
	val21 := rand.Int()
	store(t, key2, val21, storage)
	get(t, key1, []string{
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val12),
		fmt.Sprintf("%d", val13),
	}, storage)
	get(t, key2, []string{
		fmt.Sprintf("%d", val21),
	}, storage)
	val22 := rand.Int()
	store(t, key2, val22, storage)
	get(t, key1, []string{
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val12),
		fmt.Sprintf("%d", val13),
	}, storage)
	get(t, key2, []string{
		fmt.Sprintf("%d", val21),
		fmt.Sprintf("%d", val22),
	}, storage)
	val14 := rand.Int()
	store(t, key1, val14, storage)
	get(t, key1, []string{
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val11),
		fmt.Sprintf("%d", val12),
		fmt.Sprintf("%d", val13),
		fmt.Sprintf("%d", val14),
	}, storage)
	get(t, key2, []string{
		fmt.Sprintf("%d", val21),
		fmt.Sprintf("%d", val22),
	}, storage)
}

func store(t *testing.T, metricName string, v int, storage *mem.Storage) {
	val1Str := fmt.Sprintf("%d", v)
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
