package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andreevym/metric-collector/internal/compressor"
)

func (m *Middleware) ResponseGzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		oldBody := r.Body
		defer oldBody.Close()
		zr, err := gzip.NewReader(oldBody)
		if err != nil {
			io.WriteString(w, err.Error()) //nolint
			return
		}
		r.Body = zr
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequestGzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error()) //nolint
			return
		}
		defer gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(compressor.GzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
