package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetrics(t *testing.T) {
	memStorage := mem.NewStorage(nil)
	serviceHandlers := NewServiceHandlers(
		memStorage,
		nil,
	)
	router := NewRouter(serviceHandlers)
	ts := httptest.NewServer(router)
	defer ts.Close()

	metrics := []storage.Metric{}
	d1 := int64(1)
	d2 := int64(2)
	metrics = append(
		metrics,
		storage.Metric{
			ID:    "a",
			MType: storage.MTypeCounter,
			Delta: &d1,
			Value: nil,
		},
		storage.Metric{
			ID:    "b",
			MType: storage.MTypeCounter,
			Delta: &d1,
			Value: nil,
		},
		storage.Metric{
			ID:    "b",
			MType: storage.MTypeCounter,
			Delta: &d2,
			Value: nil,
		},
		storage.Metric{
			ID:    "a",
			MType: storage.MTypeCounter,
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

	metric := storage.Metric{
		ID:    "a",
		MType: storage.MTypeCounter,
		Delta: &d2,
		Value: nil,
	}
	expected, err := json.Marshal(metric)
	require.NoError(t, err)
	reqBody, err = json.Marshal(storage.Metric{
		ID:    metric.ID,
		MType: metric.MType,
	})
	require.NoError(t, err)
	statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/value/", bytes.NewBuffer(reqBody))
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, ValueMetricContentType, contentType)
	require.JSONEq(t, string(expected), get)

	metric = storage.Metric{
		ID:    "b",
		MType: storage.MTypeCounter,
		Delta: &d2,
		Value: nil,
	}
	expected, err = json.Marshal(metric)
	require.NoError(t, err)
	reqBody, err = json.Marshal(storage.Metric{
		ID:    metric.ID,
		MType: metric.MType,
	})
	require.NoError(t, err)
	statusCode, contentType, get = testRequest(t, ts, http.MethodPost, "/value/", bytes.NewBuffer(reqBody))
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, ValueMetricContentType, contentType)
	require.JSONEq(t, string(expected), get)
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
