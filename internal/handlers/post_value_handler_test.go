package handlers_test

import (
	bytes2 "bytes"
	"encoding/json"
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

func TestGetHandler(t *testing.T) {
	type want struct {
		contentType string
		resp        string
		statusCode  int
	}
	tests := []struct {
		name          string
		want          want
		request       string
		httpMethod    string
		createCounter map[string]string
		createGauge   map[string]string
		updateCounter map[string]string
		updateGauge   map[string]string
		metrics       *handlers.Metrics
	}{
		{
			name: "success get counter",
			want: want{
				contentType: handlers.ValueMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "1",
			},
			createCounter: map[string]string{
				"test":  "1",
				"test2": "2",
			},
			createGauge: map[string]string{
				"test":  "3",
				"test2": "4",
			},
			request: "/value/",
			metrics: &handlers.Metrics{
				ID:    "test",
				MType: multistorage.MetricTypeCounter,
			},
			httpMethod: http.MethodPost,
		},
		{
			name: "success get gauge",
			want: want{
				contentType: handlers.ValueMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "3",
			},
			createCounter: map[string]string{
				"test":  "1",
				"test2": "2",
			},
			createGauge: map[string]string{
				"test":  "3",
				"test2": "4",
			},
			request: "/value/",
			metrics: &handlers.Metrics{
				ID:    "test",
				MType: multistorage.MetricTypeGauge,
			},
			httpMethod: http.MethodPost,
		},
		{
			name: "success get counter after update",
			want: want{
				contentType: handlers.ValueMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "3",
			},
			createCounter: map[string]string{
				"test":  "1",
				"test2": "2",
			},
			updateCounter: map[string]string{
				"test":  "3",
				"test2": "4",
			},
			request:    "/value/",
			httpMethod: http.MethodPost,
			metrics: &handlers.Metrics{
				ID:    "test",
				MType: multistorage.MetricTypeCounter,
			},
		},
		{
			name: "success get gauge after update",
			want: want{
				contentType: handlers.ValueMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "3",
			},
			createGauge: map[string]string{
				"test":  "1",
				"test2": "2",
			},
			updateGauge: map[string]string{
				"test":  "3",
				"test2": "4",
			},
			request:    "/value/",
			httpMethod: http.MethodPost,
			metrics: &handlers.Metrics{
				ID:    "test",
				MType: multistorage.MetricTypeGauge,
			},
		},
		{
			name: "unknown 'metricName' get gauge",
			want: want{
				contentType: handlers.ValueMetricContentType,
				statusCode:  http.StatusNotFound,
				resp:        "",
			},
			request:    "/value/",
			httpMethod: http.MethodPost,
			metrics: &handlers.Metrics{
				ID:    "test",
				MType: multistorage.MetricTypeGauge,
			},
		},
		{
			name: "unknown 'metricName' get counter",
			want: want{
				contentType: handlers.ValueMetricContentType,
				statusCode:  http.StatusNotFound,
				resp:        "",
			},
			request:    "/value/",
			httpMethod: http.MethodPost,
			metrics: &handlers.Metrics{
				ID:    "test",
				MType: multistorage.MetricTypeCounter,
			},
		},
		{
			name: "unknown 'metricType'",
			want: want{
				contentType: handlers.ValueMetricContentType,
				statusCode:  http.StatusNotFound,
				resp:        "",
			},
			request:    "/value/",
			httpMethod: http.MethodPost,
			metrics: &handlers.Metrics{
				ID:    "test",
				MType: "TestGauge",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			counterMemStorage := mem.NewStorage()
			for k, v := range test.createCounter {
				err := counterMemStorage.Create(k, v)
				assert.NoError(t, err)
			}
			gaugeMemStorage := mem.NewStorage()
			for k, v := range test.createGauge {
				err := gaugeMemStorage.Create(k, v)
				assert.NoError(t, err)
			}
			for k, v := range test.updateCounter {
				err := counterMemStorage.Update(k, []string{v})
				assert.NoError(t, err)
			}
			for k, v := range test.updateGauge {
				err := gaugeMemStorage.Update(k, []string{v})
				assert.NoError(t, err)
			}
			store, err := multistorage.NewStorage(counterMemStorage, gaugeMemStorage, emptyServerConfig)
			require.NoError(t, err)
			serviceHandlers := handlers.NewServiceHandlers(store, nil)
			router := handlers.NewRouter(serviceHandlers)
			ts := httptest.NewServer(router)
			defer ts.Close()
			var reqBody []byte
			if test.metrics != nil {
				reqBody, err = json.Marshal(test.metrics)
				require.NoError(t, err)
			} else {
				reqBody = []byte{}
			}
			statusCode, contentType, get := testRequest(t, ts, test.httpMethod, test.request, bytes2.NewBuffer(reqBody))
			assert.Equal(t, test.want.statusCode, statusCode)

			if test.want.resp != "" {
				assert.Equal(t, test.want.contentType, contentType)

				respMetrics := handlers.Metrics{}
				err = json.Unmarshal([]byte(get), &respMetrics)
				require.NoError(t, err)

				if test.metrics != nil {
					if test.metrics.MType == multistorage.MetricTypeGauge {
						v, err := strconv.ParseFloat(test.want.resp, 64)
						require.NoError(t, err)
						test.metrics.Value = &v
					} else if test.metrics.MType == multistorage.MetricTypeCounter {
						v, err := strconv.ParseInt(test.want.resp, 10, 64)
						require.NoError(t, err)
						test.metrics.Delta = &v
					}

					bytes, err := json.Marshal(test.metrics)
					require.NoError(t, err)
					assert.JSONEq(t, string(bytes), get)
				}
			}
		})
	}
}
