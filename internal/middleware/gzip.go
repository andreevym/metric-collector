package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andreevym/metric-collector/internal/compressor"
	"github.com/andreevym/metric-collector/internal/logger"
	"go.uber.org/zap"
)

// ResponseGzipMiddleware returns an HTTP middleware that decompresses the request body
// if it is encoded with gzip. If the request is not gzip-encoded, it passes the request
// to the next handler in the chain without modification.
func (m *Middleware) ResponseGzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		oldBody := r.Body
		defer func(oldBody io.ReadCloser) {
			err := oldBody.Close()
			if err != nil {
				logger.Logger().Error("oldBody.Close", zap.Error(err))
			}
		}(oldBody)
		zr, err := gzip.NewReader(oldBody)
		if err != nil {
			_, err = io.WriteString(w, err.Error())
			if err != nil {
				logger.Logger().Error("value can't be written", zap.Error(err))
				return
			} //nolint
			return
		}
		r.Body = zr
		next.ServeHTTP(w, r)
	})
}

// RequestGzipMiddleware returns an HTTP middleware that compresses the response body
// with gzip if the client accepts gzip encoding. If the client does not accept gzip encoding,
// it passes the response to the next handler in the chain without modification.
func (m *Middleware) RequestGzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, err = io.WriteString(w, err.Error())
			if err != nil {
				logger.Logger().Error("value can't be written", zap.Error(err))
				return
			} //nolint
			return
		}
		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				logger.Logger().Error("gz.Close", zap.Error(err))
			}
		}(gz)
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(compressor.GzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
