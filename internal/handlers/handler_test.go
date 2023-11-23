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
	counterMemStorage := mem.NewStorage(nil)
	gaugeMemStorage := mem.NewStorage(nil)
	store, err := multistorage.NewMetricManager(counterMemStorage, gaugeMemStorage)
	require.NoError(t, err)
	serviceHandlers := handlers.NewServiceHandlers(store, nil)
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
		assert.Equal(t, string(bytes), get)
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
		assert.Equal(t, string(bytes), get)
		bytes, err = json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeGauge,
		})
		require.NoError(t, err)
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
	counterMemStorage := mem.NewStorage(nil)
	gaugeMemStorage := mem.NewStorage(nil)
	metricStorage, err := multistorage.NewMetricManager(counterMemStorage, gaugeMemStorage)
	require.NoError(t, err)
	serviceHandlers := handlers.NewServiceHandlers(metricStorage, nil)
	router := handlers.NewRouter(serviceHandlers)
	ts := httptest.NewServer(router)
	defer ts.Close()

	count := 3
	for i := 0; i < count; i++ {
		key := rand.Int()
		val1 := rand.Int63()
		bytes, err := json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeCounter,
			Delta: &val1,
		})
		require.NoError(t, err)
		statusCode, contentType, get := testRequest(t, ts, http.MethodPost, "/update/", bytes2.NewBuffer(bytes))
		require.Equal(t, http.StatusOK, statusCode)
		require.Equal(t, handlers.UpdateMetricContentType, contentType)
		require.Equal(t, string(bytes), get)
		val2 := rand.Int63()
		bytes, err = json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeCounter,
			Delta: &val2,
		})
		require.NoError(t, err)

		res := val1 + val2
		resBytes, err := json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeCounter,
			Delta: &res,
		})
		require.NoError(t, err)

		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/update/", bytes2.NewBuffer(bytes))
		require.Equal(t, http.StatusOK, statusCode)
		require.Equal(t, handlers.UpdateMetricContentType, contentType)
		require.Equal(t, string(resBytes), get)
		bytes, err = json.Marshal(handlers.Metrics{
			ID:    strconv.Itoa(key),
			MType: multistorage.MetricTypeCounter,
		})
		require.NoError(t, err)
		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/value/", bytes2.NewBuffer(bytes))
		require.Equal(t, http.StatusOK, statusCode)
		require.Equal(t, handlers.ValueMetricContentType, contentType)
		require.JSONEq(t, string(resBytes), get)
	}
}
