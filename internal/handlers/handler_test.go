package handlers_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
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
		val1 := rand.Int()
		statusCode, contentType, get := testRequest(t, ts, http.MethodPost, fmt.Sprintf("/update/gauge/test%d/%d", key, val1))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.UpdateMetricContentType, contentType)
		assert.Equal(t, "", get)
		val2 := rand.Int()
		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, fmt.Sprintf("/update/gauge/test%d/%d", key, val2))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.UpdateMetricContentType, contentType)
		assert.Equal(t, "", get)
		statusCode, contentType, get = testRequest(t, ts, http.MethodGet, fmt.Sprintf("/value/gauge/test%d", key))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, "text/plain; charset=utf-8", contentType)
		assert.Equal(t, fmt.Sprintf("%d", val2), get)
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
		val1 := rand.Int()
		statusCode, contentType, get := testRequest(t, ts, http.MethodPost, fmt.Sprintf("/update/counter/test%d/%d", key, val1))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.UpdateMetricContentType, contentType)
		assert.Equal(t, "", get)
		val2 := rand.Int()
		statusCode, contentType, get = testRequest(t, ts, http.MethodPost, fmt.Sprintf("/update/counter/test%d/%d", key, val2))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, handlers.UpdateMetricContentType, contentType)
		assert.Equal(t, "", get)
		statusCode, contentType, get = testRequest(t, ts, http.MethodGet, fmt.Sprintf("/value/counter/test%d", key))
		assert.Equal(t, http.StatusOK, statusCode)
		assert.Equal(t, "text/plain; charset=utf-8", contentType)
		assert.Equal(t, fmt.Sprintf("%d", val1+val2), get)
	}
}
