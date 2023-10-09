package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/assert"
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
	}{
		{
			name: "success get counter",
			want: want{
				contentType: "text/plain; charset=utf-8",
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
			request:    "/value/counter/test",
			httpMethod: http.MethodGet,
		},
		{
			name: "success get gauge",
			want: want{
				contentType: "text/plain; charset=utf-8",
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
			request:    "/value/gauge/test",
			httpMethod: http.MethodGet,
		},
		{
			name: "success get counter after update",
			want: want{
				contentType: "text/plain; charset=utf-8",
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
			request:    "/value/counter/test",
			httpMethod: http.MethodGet,
		},
		{
			name: "success get gauge after update",
			want: want{
				contentType: "text/plain; charset=utf-8",
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
			request:    "/value/gauge/test",
			httpMethod: http.MethodGet,
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
			server := handlers.NewServer(counterMemStorage, gaugeMemStorage)
			ts := httptest.NewServer(handlers.Router(server))
			defer ts.Close()

			statusCode, contentType, get := testRequest(t, ts, test.httpMethod, test.request)
			assert.Equal(t, test.want.statusCode, statusCode)
			assert.Equal(t, test.want.contentType, contentType)
			assert.Equal(t, test.want.resp, get)
		})
	}
}
