package handlers_test

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GaugeEndToEnd(t *testing.T) {
	memStorage := mem.NewStorage(nil)
	serviceHandlers := handlers.NewServiceHandlers(memStorage, nil)
	router := handlers.NewRouter(serviceHandlers)
	ts := httptest.NewServer(router)
	defer ts.Close()

	count := 3
	for i := 0; i < count; i++ {
		key := rand.Int()
		val1 := rand.Float64()
		bytes, err := json.Marshal(storage.Metric{
			ID:    strconv.Itoa(key),
			MType: storage.MTypeGauge,
			Value: &val1,
		})
		require.NoError(t, err)
		statusCode, contentType, get := testRequest(t, ts, http.MethodPost, "/update/", bytes)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.UpdateMetricContentType, contentType)
		assert.Equal(t, string(bytes), get)
		val2 := rand.Float64()
		bytes, err = json.Marshal(storage.Metric{
			ID:    strconv.Itoa(key),
			MType: storage.MTypeGauge,
			Value: &val2,
		})
		require.NoError(t, err)
		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/update/", bytes)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.UpdateMetricContentType, contentType)
		assert.Equal(t, string(bytes), get)
		bytes, err = json.Marshal(storage.Metric{
			ID:    strconv.Itoa(key),
			MType: storage.MTypeGauge,
		})
		require.NoError(t, err)
		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/value/", bytes)
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.ValueMetricContentType, contentType)
		bytes, err = json.Marshal(storage.Metric{
			ID:    strconv.Itoa(key),
			MType: storage.MTypeGauge,
			Value: &val2,
		})
		require.NoError(t, err)
		assert.JSONEq(t, string(bytes), get)
	}
}

func TestHandler_CounterEndToEnd(t *testing.T) {
	memStorage := mem.NewStorage(nil)
	serviceHandlers := handlers.NewServiceHandlers(memStorage, nil)
	router := handlers.NewRouter(serviceHandlers)
	ts := httptest.NewServer(router)
	defer ts.Close()

	count := 3
	for i := 0; i < count; i++ {
		key := rand.Int()
		val1 := rand.Int63()
		bytes, err := json.Marshal(storage.Metric{
			ID:    strconv.Itoa(key),
			MType: storage.MTypeCounter,
			Delta: &val1,
		})
		require.NoError(t, err)
		statusCode, contentType, get := testRequest(t, ts, http.MethodPost, "/update/", bytes)
		require.Equal(t, http.StatusOK, statusCode)
		require.Equal(t, handlers.UpdateMetricContentType, contentType)
		require.Equal(t, string(bytes), get)
		val2 := rand.Int63()
		bytes, err = json.Marshal(storage.Metric{
			ID:    strconv.Itoa(key),
			MType: storage.MTypeCounter,
			Delta: &val2,
		})
		require.NoError(t, err)

		res := val1 + val2
		resBytes, err := json.Marshal(storage.Metric{
			ID:    strconv.Itoa(key),
			MType: storage.MTypeCounter,
			Delta: &res,
		})
		require.NoError(t, err)

		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/update/", bytes)
		require.Equal(t, http.StatusOK, statusCode)
		require.Equal(t, handlers.UpdateMetricContentType, contentType)
		require.Equal(t, string(resBytes), get)
		bytes, err = json.Marshal(storage.Metric{
			ID:    strconv.Itoa(key),
			MType: storage.MTypeCounter,
		})
		require.NoError(t, err)
		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/value/", bytes)
		require.Equal(t, http.StatusOK, statusCode)
		require.Equal(t, handlers.ValueMetricContentType, contentType)
		require.JSONEq(t, string(resBytes), get)
	}
}
