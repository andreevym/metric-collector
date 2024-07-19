package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/andreevym/metric-collector/internal/crypto"
)

func (m *Middleware) RequestCryptoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.CryptoKey == "" {
			next.ServeHTTP(w, r)
			return
		}
		readAll, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		decode, err := crypto.Decode(m.CryptoKey, string(readAll))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader([]byte(decode)))
		next.ServeHTTP(w, r)
	})
}
