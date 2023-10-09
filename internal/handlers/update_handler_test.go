package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreevym/metric-collector/internal/handlers"
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
				contentType: "text/plain",
				statusCode:  http.StatusOK,
				resp:        "",
			},
			request:    "/update/counter/test/1",
			httpMethod: http.MethodPost,
		},
		{
			name: "success update gauge",
			want: want{
				contentType: "text/plain",
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

			request := httptest.NewRequest(test.httpMethod, test.request, nil)
			w := httptest.NewRecorder()
			h := handlers.UpdateHandler(counterMemStorage, gaugeMemStorage)
			h(w, request)

			result := w.Result()

			assert.Equal(t, test.want.statusCode, result.StatusCode)
			assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))

			resBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			require.Empty(t, resBody)
		})
	}
}
