package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreevym/metric-collector/internal/compressor"
	"github.com/andreevym/metric-collector/internal/middleware"
	"github.com/andreevym/metric-collector/internal/multistorage"
	"github.com/andreevym/metric-collector/internal/storage/mem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ...

func TestGzipCompressionUpdate(t *testing.T) {
	counterMemStorage := mem.NewStorage()
	err := counterMemStorage.Create("A", "1")
	assert.NoError(t, err)

	gaugeMemStorage := mem.NewStorage()
	err = gaugeMemStorage.Create("B", "0.2")
	assert.NoError(t, err)

	store, err := multistorage.NewStorage(counterMemStorage, gaugeMemStorage)
	require.NoError(t, err)
	serviceHandlers := NewServiceHandlers(store)
	router := NewRouter(serviceHandlers)
	srv := httptest.NewServer(middleware.GzipMiddleware(router))
	defer srv.Close()

	t.Run("without_gzip", func(t *testing.T) {
		valueA := int64(2)
		requestBody, err := json.Marshal(&Metrics{
			ID:    "N",
			MType: "counter",
			Delta: &valueA,
		})
		require.NoError(t, err)

		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&Metrics{
			ID:    "N",
			MType: "counter",
			Delta: &valueA,
		})
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Accept-Encoding", "")
		_, _, respBody := testCompressRequest(t, srv, http.MethodPost, "/update", bytes.NewBuffer(requestBody), header)

		require.NoError(t, err)

		require.JSONEq(t, string(successBody), respBody)
	})

	t.Run("sends_gzip", func(t *testing.T) {
		valueA := int64(2)
		requestBody, err := json.Marshal(&Metrics{
			ID:    "A",
			MType: "counter",
			Delta: &valueA,
		})
		require.NoError(t, err)

		newValueA := int64(3)
		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&Metrics{
			ID:    "A",
			MType: "counter",
			Delta: &newValueA,
		})
		require.NoError(t, err)

		compressed, err := compressor.Compress(requestBody)
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Content-Encoding", "gzip")
		_, _, respBody := testCompressRequest(t, srv, http.MethodPost, "/update", bytes.NewBuffer(compressed), header)

		decompressed, err := compressor.Decompress([]byte(respBody))
		require.NoError(t, err)

		require.JSONEq(t, string(successBody), string(decompressed))
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		valueA := int64(2)
		requestBody, err := json.Marshal(&Metrics{
			ID:    "B",
			MType: "counter",
			Delta: &valueA,
		})
		require.NoError(t, err)

		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&Metrics{
			ID:    "B",
			MType: "counter",
			Delta: &valueA,
		})
		require.NoError(t, err)

		compressed, err := compressor.Compress(requestBody)
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Content-Encoding", "gzip")
		header.Set("Accept-Encoding", "gzip")
		_, _, respBody := testCompressRequest(t, srv, http.MethodPost, "/update", bytes.NewBuffer(compressed), header)

		decompressed, err := compressor.Decompress([]byte(respBody))
		require.NoError(t, err)

		require.JSONEq(t, string(successBody), string(decompressed))
	})
}

func TestGzipCompressionValue(t *testing.T) {
	counterMemStorage := mem.NewStorage()
	err := counterMemStorage.Create("A", "1")
	assert.NoError(t, err)

	gaugeMemStorage := mem.NewStorage()
	err = gaugeMemStorage.Create("B", "0.2")
	assert.NoError(t, err)

	store, err := multistorage.NewStorage(counterMemStorage, gaugeMemStorage)
	require.NoError(t, err)
	serviceHandlers := NewServiceHandlers(store)
	router := NewRouter(serviceHandlers)
	srv := httptest.NewServer(middleware.GzipMiddleware(router))
	defer srv.Close()

	t.Run("without_gzip", func(t *testing.T) {
		requestBody, err := json.Marshal(&Metrics{
			ID:    "A",
			MType: "counter",
		})
		require.NoError(t, err)

		valueA := int64(1)
		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&Metrics{
			ID:    "A",
			MType: "counter",
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
		requestBody, err := json.Marshal(&Metrics{
			ID:    "A",
			MType: "counter",
		})
		require.NoError(t, err)

		valueA := int64(1)
		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&Metrics{
			ID:    "A",
			MType: "counter",
			Delta: &valueA,
		})
		require.NoError(t, err)

		compressed, err := compressor.Compress(requestBody)
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Content-Encoding", "gzip")
		_, _, respBody := testCompressRequest(t, srv, http.MethodPost, "/value", bytes.NewBuffer(compressed), header)

		decompressed, err := compressor.Decompress([]byte(respBody))
		require.NoError(t, err)

		require.JSONEq(t, string(successBody), string(decompressed))
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		requestBody, err := json.Marshal(&Metrics{
			ID:    "A",
			MType: "counter",
		})
		require.NoError(t, err)

		valueA := int64(1)
		// ожидаемое содержимое тела ответа при успешном запросе
		successBody, err := json.Marshal(&Metrics{
			ID:    "A",
			MType: "counter",
			Delta: &valueA,
		})
		require.NoError(t, err)

		compressed, err := compressor.Compress(requestBody)
		require.NoError(t, err)

		header := http.Header{}
		header.Set("Content-Encoding", "gzip")
		header.Set("Accept-Encoding", "gzip")
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
