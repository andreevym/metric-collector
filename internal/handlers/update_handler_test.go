package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/hash"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateHandler(t *testing.T) {
	type want struct {
		contentType string
		resp        string
		statusCode  int
	}
	f := float64(1)
	d := int64(1)
	tests := []struct {
		name       string
		want       want
		request    string
		metrics    *storage.Metric
		httpMethod string
	}{
		{
			name: "success update counter",
			want: want{
				contentType: handlers.UpdateMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "{\"id\":\"test\",\"type\":\"counter\",\"delta\":1}",
			},
			request:    "/update",
			httpMethod: http.MethodPost,
			metrics: &storage.Metric{
				ID:    "test",
				MType: storage.MTypeCounter,
				Delta: &d,
			},
		},
		{
			name: "success update counter",
			want: want{
				contentType: handlers.UpdateMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "{\"id\":\"test\",\"type\":\"counter\",\"delta\":1}",
			},
			request:    "/update/",
			httpMethod: http.MethodPost,
			metrics: &storage.Metric{
				ID:    "test",
				MType: storage.MTypeCounter,
				Delta: &d,
			},
		},
		{
			name: "success update gauge",
			want: want{
				contentType: handlers.UpdateMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "{\"id\":\"test\",\"type\":\"gauge\",\"value\":1}",
			},
			request:    "/update",
			httpMethod: http.MethodPost,
			metrics: &storage.Metric{
				ID:    "test",
				MType: storage.MTypeGauge,
				Value: &f,
			},
		},
		{
			name: "success update gauge",
			want: want{
				contentType: handlers.UpdateMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "{\"id\":\"test\",\"type\":\"gauge\",\"value\":1}",
			},
			request:    "/update/",
			httpMethod: http.MethodPost,
			metrics: &storage.Metric{
				ID:    "test",
				MType: storage.MTypeGauge,
				Value: &f,
			},
		},
		{
			name: "success update counter",
			want: want{
				contentType: handlers.UpdateMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "{\"id\":\"test\",\"type\":\"counter\",\"delta\":1}",
			},
			request:    "/update/counter/test/1",
			httpMethod: http.MethodPost,
		},
		{
			name: "success update gauge",
			want: want{
				contentType: handlers.UpdateMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "{\"id\":\"test\",\"type\":\"gauge\",\"value\":1}",
			},
			request:    "/update/gauge/test/1",
			httpMethod: http.MethodPost,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			memStorage := mem.NewStorage(nil)
			serviceHandlers := handlers.NewServiceHandlers(memStorage, nil)
			router := handlers.NewRouter(serviceHandlers)
			ts := httptest.NewServer(router)
			defer ts.Close()
			var bytes []byte
			if test.metrics != nil {
				marshal, err := json.Marshal(test.metrics)
				require.NoError(t, err)
				bytes = marshal
			}
			statusCode, contentType, get := testRequest(t, ts, test.httpMethod, test.request, bytes)
			assert.Equal(t, test.want.statusCode, statusCode)
			assert.Equal(t, test.want.contentType, contentType)
			assert.Equal(t, test.want.resp, get)
		})
	}
}

func signedTestRequest(
	t *testing.T,
	ts *httptest.Server,
	method, path string,
	reqBody []byte,
	secretKey string,
) (int, string, string) {
	headerMap := make(map[string]string)
	if secretKey != "" && len(reqBody) > 0 {
		encodeHash := hash.EncodeHash(reqBody, secretKey)
		headerMap[middleware.HashHeaderKey] = encodeHash
	}

	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	for k, v := range headerMap {
		req.Header.Add(k, v)
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = resp.Body.Close()
	require.NoError(t, err)

	contentType := resp.Header.Get("Content-Type")
	return resp.StatusCode, contentType, string(respBody)
}

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method string,
	path string,
	reqBody []byte,
) (int, string, string) {
	return signedTestRequest(t, ts, method, path, reqBody, "")
}

func TestUpdates(t *testing.T) {
	idCounter := "CounterBatchZip" + strconv.Itoa(rand.Int())
	idGauge := "GaugeBatchZip" + strconv.Itoa(rand.Int())
	valueCounter1, valueCounter2 := int64(rand.Int()), int64(rand.Int())
	valueGauge1, valueGauge2 := float64(rand.Float32()), float64(rand.Float32())

	metrics := []storage.Metric{
		{
			ID:    idCounter,
			MType: "counter",
			Delta: &valueCounter1,
		},
		{
			ID:    idGauge,
			MType: "gauge",
			Value: &valueGauge1,
		},
		{
			ID:    idCounter,
			MType: "counter",
			Delta: &valueCounter2,
		},
		{
			ID:    idGauge,
			MType: "gauge",
			Value: &valueGauge2,
		},
	}

	bytes, err := json.Marshal(metrics)
	require.NoError(t, err)
	memStorage := mem.NewStorage(nil)
	serviceHandlers := handlers.NewServiceHandlers(memStorage, nil)
	router := handlers.NewRouter(serviceHandlers)
	ts := httptest.NewServer(router)
	defer ts.Close()
	statusCode, contentType, get := testRequest(t, ts, http.MethodPost, "/updates/", bytes)
	assert.Equal(t, statusCode, statusCode)
	assert.Equal(t, contentType, contentType)
	assert.Equal(t, "", get)

	statusCode, contentType, get = testRequest(t, ts, http.MethodGet, "/value/counter/"+idCounter, nil)
	assert.Equal(t, statusCode, statusCode)
	assert.Equal(t, contentType, contentType)
	assert.Equal(t, strconv.FormatInt(valueCounter1+valueCounter2, 10), get)

	statusCode, contentType, get = testRequest(t, ts, http.MethodGet, "/value/gauge/"+idGauge, nil)
	assert.Equal(t, statusCode, statusCode)
	assert.Equal(t, contentType, contentType)
	res := strconv.FormatFloat(valueGauge2, 'f', -1, 64)
	assert.Equal(t, res, get)
}
