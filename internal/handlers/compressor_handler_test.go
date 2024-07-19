package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreevym/metric-collector/internal/compressor"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGzipCompressionUpdate(t *testing.T) {
	memStorage := mem.NewStorage(nil)
	i := int64(1)
	err := memStorage.Create(context.TODO(), &store.Metric{
		ID:    "B",
		MType: store.MTypeCounter,
		Delta: &i,
	})
	assert.NoError(t, err)
	err = memStorage.Create(context.TODO(), &store.Metric{
		ID:    "A",
		MType: store.MTypeCounter,
		Delta: &i,
	})
	assert.NoError(t, err)
	f := 0.2
	err = memStorage.Create(context.TODO(), &store.Metric{
		ID:    "B",
		MType: store.MTypeGauge,
		Value: &f,
	})
	assert.NoError(t, err)

	serviceHandlers := NewServiceHandlers(memStorage, nil)
	m := middleware.NewMiddleware("", "")
	router := NewRouter(
		serviceHandlers,
		m.RequestGzipMiddleware,
		m.ResponseGzipMiddleware,
	)
	srv := httptest.NewServer(router)
	defer srv.Close()

	t.Run("without_gzip", func(t *testing.T) {
		valueA := int64(2)
		requestBody, err := json.Marshal(&store.Metric{
			ID:    "N",
			MType: store.MTypeCounter,
			Delta: &valueA,
		})
		require.NoError(t, err)

		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&store.Metric{
			ID:    "N",
			MType: store.MTypeCounter,
			Delta: &valueA,
		})
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Accept-Encoding", "")
		_, _, respBody := testCompressRequest(t, srv, http.MethodPost, "/update", bytes.NewBuffer(requestBody), header)

		require.JSONEq(t, string(successBody), respBody)
	})

	t.Run("sends_gzip", func(t *testing.T) {
		valueA := int64(2)
		requestBody, err := json.Marshal(&store.Metric{
			ID:    "A",
			MType: store.MTypeCounter,
			Delta: &valueA,
		})
		require.NoError(t, err)

		newValueA := int64(3)
		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&store.Metric{
			ID:    "A",
			MType: store.MTypeCounter,
			Delta: &newValueA,
		})
		require.NoError(t, err)

		compressed, err := compressor.Compress(requestBody)
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Content-Encoding", "gzip")
		header.Set("Content-Type", "application/json")
		_, _, respBody := testCompressRequest(t, srv, http.MethodPost, "/update", bytes.NewBuffer(compressed), header)

		require.JSONEq(t, string(successBody), respBody)
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		valueRequest := int64(2)
		requestBody, err := json.Marshal(&store.Metric{
			ID:    "B",
			MType: store.MTypeCounter,
			Delta: &valueRequest,
		})
		require.NoError(t, err)

		// ожидаемое содержимое тела ответа при успешном запросе
		valueAfterUpdate := int64(3)
		successBody, err := json.Marshal(&store.Metric{
			ID:    "B",
			MType: store.MTypeCounter,
			Delta: &valueAfterUpdate,
		})
		require.NoError(t, err)

		compressed, err := compressor.Compress(requestBody)
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Content-Encoding", "gzip")
		header.Set("Accept-Encoding", "gzip")
		header.Set("Content-Type", "application/json")
		_, _, respBody := testCompressRequest(t, srv, http.MethodPost, "/update", bytes.NewBuffer(compressed), header)

		decompressed, err := compressor.Decompress([]byte(respBody))
		require.NoError(t, err)

		require.JSONEq(t, string(successBody), string(decompressed))
	})
}

func TestGzipCompressionValue(t *testing.T) {
	memStorage := mem.NewStorage(nil)
	valueA := int64(1)
	err := memStorage.Create(context.TODO(), &store.Metric{
		ID:    "A",
		MType: store.MTypeCounter,
		Delta: &valueA,
	})
	assert.NoError(t, err)

	serviceHandlers := NewServiceHandlers(memStorage, nil)
	m := middleware.NewMiddleware("", "")
	router := NewRouter(
		serviceHandlers,
		m.ResponseGzipMiddleware,
		m.RequestGzipMiddleware,
	)
	srv := httptest.NewServer(router)
	defer srv.Close()

	t.Run("without_gzip", func(t *testing.T) {
		requestBody, err := json.Marshal(&store.Metric{
			ID:    "A",
			MType: store.MTypeCounter,
		})
		require.NoError(t, err)

		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&store.Metric{
			ID:    "A",
			MType: store.MTypeCounter,
			Delta: &valueA,
		})
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Accept-Encoding", "")
		_, _, respBody := testCompressRequest(t, srv, http.MethodPost, "/value", bytes.NewBuffer(requestBody), header)

		require.NoError(t, err)

		require.JSONEq(t, string(successBody), respBody)
	})

	t.Run("sends_gzip", func(t *testing.T) {
		requestBody, err := json.Marshal(&store.Metric{
			ID:    "A",
			MType: store.MTypeCounter,
		})
		require.NoError(t, err)

		valueA := int64(1)
		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&store.Metric{
			ID:    "A",
			MType: store.MTypeCounter,
			Delta: &valueA,
		})
		require.NoError(t, err)

		compressed, err := compressor.Compress(requestBody)
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Content-Encoding", "gzip")
		header.Set("Content-Type", "application/json")
		_, _, respBody := testCompressRequest(t, srv, http.MethodPost, "/value", bytes.NewBuffer(compressed), header)

		require.JSONEq(t, string(successBody), respBody)
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		requestBody, err := json.Marshal(&store.Metric{
			ID:    "A",
			MType: store.MTypeCounter,
		})
		require.NoError(t, err)

		valueA := int64(1)
		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&store.Metric{
			ID:    "A",
			MType: store.MTypeCounter,
			Delta: &valueA,
		})
		require.NoError(t, err)

		compressed, err := compressor.Compress(requestBody)
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Content-Encoding", "gzip")
		header.Set("Accept-Encoding", "gzip")
		header.Set("Content-Type", "application/json")
		_, _, respBody := testCompressRequest(t, srv, http.MethodPost, "/value", bytes.NewBuffer(compressed), header)

		decompressed, err := compressor.Decompress([]byte(respBody))
		require.NoError(t, err)

		require.JSONEq(t, string(successBody), string(decompressed))
	})
}

func testCompressRequest(
	t *testing.T,
	ts *httptest.Server,
	method, path string,
	reqBody io.Reader,
	header http.Header,
) (int, string, string) {
	req, err := http.NewRequest(method, ts.URL+path, reqBody)
	require.NoError(t, err)

	req.Header = header

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = resp.Body.Close()
	require.NoError(t, err)

	contentType := resp.Header.Get("Content-Type")
	return resp.StatusCode, contentType, string(respBody)
}
