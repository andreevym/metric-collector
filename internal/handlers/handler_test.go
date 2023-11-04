package handlers_test

import (
	bytes2 "bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GaugeEndToEnd(t *testing.T) {
	counterMemStorage := mem.NewStorage()
	gaugeMemStorage := mem.NewStorage()
	store, err := multistorage.NewStorage(counterMemStorage, gaugeMemStorage)
	require.NoError(t, err)
	serviceHandlers := handlers.NewServiceHandlers(store)
	router := handlers.NewRouter(serviceHandlers)
	ts := httptest.NewServer(router)
	defer ts.Close()

	count := 3
	for i := 0; i < count; i++ {
		key := rand.Int()
		val1 := rand.Float64()
		bytes, err := json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeGauge,
			Value: &val1,
		})
		require.NoError(t, err)
		statusCode, contentType, get := testRequest(t, ts, http.MethodPost, "/update/", bytes2.NewBuffer(bytes))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.UpdateMetricContentType, contentType)
		assert.Equal(t, "", get)
		val2 := rand.Float64()
		bytes, err = json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeGauge,
			Value: &val2,
		})
		require.NoError(t, err)
		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/update/", bytes2.NewBuffer(bytes))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.UpdateMetricContentType, contentType)
		assert.Equal(t, "", get)
		bytes, err = json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeGauge,
		})
		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/value/", bytes2.NewBuffer(bytes))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.ValueMetricContentType, contentType)
		bytes, err = json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeGauge,
			Value: &val2,
		})
		require.NoError(t, err)
		assert.JSONEq(t, string(bytes), get)
	}
}

func TestHandler_CounterEndToEnd(t *testing.T) {
	counterMemStorage := mem.NewStorage()
	gaugeMemStorage := mem.NewStorage()
	store, err := multistorage.NewStorage(counterMemStorage, gaugeMemStorage)
	require.NoError(t, err)
	serviceHandlers := handlers.NewServiceHandlers(store)
	router := handlers.NewRouter(serviceHandlers)
	ts := httptest.NewServer(router)
	defer ts.Close()

	count := 3
	for i := 0; i < count; i++ {
		key := rand.Int()
		val1 := float64(rand.Int())
		bytes, err := json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeCounter,
			Value: &val1,
		})
		require.NoError(t, err)
		statusCode, contentType, get := testRequest(t, ts, http.MethodPost, "/update/", bytes2.NewBuffer(bytes))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.UpdateMetricContentType, contentType)
		assert.Equal(t, "", get)
		val2 := float64(rand.Int())
		bytes, err = json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeCounter,
			Value: &val2,
		})
		require.NoError(t, err)
		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/update/", bytes2.NewBuffer(bytes))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.UpdateMetricContentType, contentType)
		assert.Equal(t, "", get)
		bytes, err = json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeCounter,
		})
		require.NoError(t, err)
		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/value/", bytes2.NewBuffer(bytes))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.ValueMetricContentType, contentType)
		f := val1 + val2
		bytes, err = json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeCounter,
			Value: &f,
		})
		require.NoError(t, err)
		assert.JSONEq(t, string(bytes), get)
	}
}
