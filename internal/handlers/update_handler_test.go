package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/multistorage"
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
	tests := []struct {
		name       string
		want       want
		request    string
		httpMethod string
	}{
		{
			name: "success update counter",
			want: want{
				contentType: handlers.UpdateMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "",
			},
			request:    "/update/counter/test/1",
			httpMethod: http.MethodPost,
		},
		{
			name: "success update gauge",
			want: want{
				contentType: handlers.UpdateMetricContentType,
				statusCode:  http.StatusOK,
				resp:        "",
			},
			request:    "/update/gauge/test/1",
			httpMethod: http.MethodPost,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			counterMemStorage := mem.NewStorage()
			gaugeMemStorage := mem.NewStorage()
			store, err := multistorage.NewStorage(counterMemStorage, gaugeMemStorage)
			require.NoError(t, err)
			serviceHandlers := handlers.NewServiceHandlers(store)
			router := handlers.NewRouter(serviceHandlers)
			ts := httptest.NewServer(router)
			defer ts.Close()

			statusCode, contentType, get := testRequest(t, ts, test.httpMethod, test.request)
			assert.Equal(t, test.want.statusCode, statusCode)
			assert.Equal(t, test.want.contentType, contentType)
			assert.Equal(t, test.want.resp, get)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
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
