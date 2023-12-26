package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/storage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const secretKey = "testpassword"

func TestPostHandler(t *testing.T) {
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
		metrics       *storage.Metric
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
			metrics: &storage.Metric{
				ID:    "test",
				MType: storage.MTypeCounter,
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
			metrics: &storage.Metric{
				ID:    "test",
				MType: storage.MTypeGauge,
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
			metrics: &storage.Metric{
				ID:    "test",
				MType: storage.MTypeCounter,
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
			metrics: &storage.Metric{
				ID:    "test",
				MType: storage.MTypeGauge,
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
			metrics: &storage.Metric{
				ID:    "test",
				MType: storage.MTypeGauge,
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
			metrics: &storage.Metric{
				ID:    "test",
				MType: storage.MTypeCounter,
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
			metrics: &storage.Metric{
				ID:    "test",
				MType: "TestGauge",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			memStorage := mem.NewStorage(nil)
			for k, v := range test.createCounter {
				i, _ := strconv.ParseInt(v, 10, 64)
				metric := &storage.Metric{
					ID:    k,
					MType: storage.MTypeCounter,
					Delta: &i,
					Value: nil,
				}
				err := memStorage.Create(context.TODO(), metric)
				assert.NoError(t, err)
			}
			for k, v := range test.createGauge {
				i, _ := strconv.ParseFloat(v, 64)
				metric := &storage.Metric{
					ID:    k,
					MType: storage.MTypeGauge,
					Delta: nil,
					Value: &i,
				}
				err := memStorage.Create(context.TODO(), metric)
				assert.NoError(t, err)
			}
			for k, v := range test.updateCounter {
				i, _ := strconv.ParseInt(v, 10, 64)
				metric := &storage.Metric{
					ID:    k,
					MType: storage.MTypeCounter,
					Delta: &i,
					Value: nil,
				}
				err := memStorage.Update(context.TODO(), metric)
				assert.NoError(t, err)
			}
			for k, v := range test.updateGauge {
				i, _ := strconv.ParseFloat(v, 64)
				metric := &storage.Metric{
					ID:    k,
					MType: storage.MTypeGauge,
					Delta: nil,
					Value: &i,
				}
				err := memStorage.Update(context.TODO(), metric)
				assert.NoError(t, err)
			}
			serviceHandlers := handlers.NewServiceHandlers(memStorage, nil)
			m := middleware.NewMiddleware(secretKey)
			router := handlers.NewRouter(serviceHandlers, m.RequestHashMiddleware)
			ts := httptest.NewServer(router)
			defer ts.Close()
			var reqBody []byte
			var err error
			if test.metrics != nil {
				reqBody, err = json.Marshal(test.metrics)
				require.NoError(t, err)
			} else {
				reqBody = []byte{}
			}
			statusCode, contentType, got := testRequest(t, ts, test.httpMethod, test.request, reqBody)
			assert.Equal(t, test.want.statusCode, statusCode)

			if test.want.resp != "" {
				assert.Equal(t, test.want.contentType, contentType)

				respMetrics := storage.Metric{}
				err = json.Unmarshal([]byte(got), &respMetrics)
				require.NoError(t, err)

				if test.metrics != nil {
					if test.metrics.MType == storage.MTypeGauge {
						v, err := strconv.ParseFloat(test.want.resp, 64)
						require.NoError(t, err)
						test.metrics.Value = &v
					} else if test.metrics.MType == storage.MTypeCounter {
						v, err := strconv.ParseInt(test.want.resp, 10, 64)
						require.NoError(t, err)
						test.metrics.Delta = &v
					}

					bytes, err := json.Marshal(test.metrics)
					require.NoError(t, err)
					assert.JSONEq(t, string(bytes), got)
				}
			}
		})
	}
}
