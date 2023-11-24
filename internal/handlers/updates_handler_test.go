package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetrics(t *testing.T) {
	counterMemStorage := mem.NewStorage(nil)
	gaugeMemStorage := mem.NewStorage(nil)
	manager, err := multistorage.NewMetricManager(counterMemStorage, gaugeMemStorage)
	require.NoError(t, err)
	serviceHandlers := NewServiceHandlers(
		manager,
		nil,
	)
	router := NewRouter(serviceHandlers)
	ts := httptest.NewServer(router)
	defer ts.Close()

	metrics := []Metrics{}
	d1 := int64(1)
	d2 := int64(2)
	metrics = append(
		metrics,
		Metrics{
			ID:    "a",
			MType: multistorage.MetricTypeCounter,
			Delta: &d1,
			Value: nil,
		},
		Metrics{
			ID:    "b",
			MType: multistorage.MetricTypeCounter,
			Delta: &d1,
			Value: nil,
		},
		Metrics{
			ID:    "a",
			MType: multistorage.MetricTypeCounter,
			Delta: &d2,
			Value: nil,
		},
	)
	reqBody, err := json.Marshal(&metrics)
	require.NoError(t, err)
	statusCode, contentType, get := testRequest(t, ts, http.MethodPost, "/updates/", bytes.NewBuffer(reqBody))
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, UpdatesMetricContentType, contentType)
	require.Equal(t, "", get)

	statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/updates/", bytes.NewBuffer(reqBody))
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, UpdatesMetricContentType, contentType)
	require.Equal(t, "", get)
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, reqBody io.Reader) (int, string, string) {
	req, err := http.NewRequest(method, ts.URL+path, reqBody)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = resp.Body.Close()
	require.NoError(t, err)

	contentType := resp.Header.Get("Content-Type")
	return resp.StatusCode, contentType, string(respBody)
}
