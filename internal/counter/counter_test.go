package counter

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/require"
)

func TestCounterEndToEnd(t *testing.T) {
	memStorage := mem.NewStorage()
	key1 := rand.Int63()
	val1 := rand.Int63()
	store(t, key1, val1, memStorage)
	get(t, key1, val1, memStorage)
	store(t, key1, val1, memStorage)
	get(t, key1, val1+val1, memStorage)
	val2 := rand.Int63()
	store(t, key1, val2, memStorage)
	get(t, key1, val1+val1+val2, memStorage)
	val3 := rand.Int63()
	store(t, key1, val3, memStorage)
	get(t, key1, val1+val1+val2+val3, memStorage)
	key2 := rand.Int63()
	val21 := rand.Int63()
	store(t, key2, val21, memStorage)
	get(t, key1, val1+val1+val2+val3, memStorage)
	val22 := rand.Int63()
	store(t, key2, val22, memStorage)
	get(t, key1, val1+val1+val2+val3, memStorage)
	val14 := rand.Int63()
	store(t, key1, val14, memStorage)
	get(t, key1, val1+val1+val2+val3+val14, memStorage)
}

func store(t *testing.T, name int64, v int64, storage *mem.Storage) {
	metricName := fmt.Sprintf("a%v", name)
	val1Str := fmt.Sprintf("%v", v)
	err := Validate(val1Str)
	require.NoError(t, err)
	_, err = Store(storage, metricName, val1Str)
	require.NoError(t, err)
}

func get(t *testing.T, name int64, expectedValue int64, storage *mem.Storage) {
	metricName := fmt.Sprintf("a%v", name)
	val1Str := fmt.Sprintf("%v", expectedValue)
	get, err := Get(storage, metricName)
	require.NoError(t, err)
	require.Equal(t, val1Str, get)
}
