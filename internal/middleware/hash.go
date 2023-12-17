package middleware

import (
	"io"
	"net/http"

	"github.com/andreevym/metric-collector/internal/hash"
	"github.com/andreevym/metric-collector/internal/logger"
	"go.uber.org/zap"
)

const hashHeaderKey = "HashSHA256"

func (m *Middleware) RequestHashMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		agentRequestBodyHash := r.Header.Get(hashHeaderKey)
		if agentRequestBodyHash != "" {
			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				logger.Logger().Error(
					"could not read request body when agent hash sha 256 is defined",
					zap.String("agentRequestBodyHash", agentRequestBodyHash),
					zap.Error(err),
				)
				return
			}
			serverRequestBodyHash := hash.Hash(bytes, m.SecretKey)
			if serverRequestBodyHash != agentRequestBodyHash {
				w.WriteHeader(http.StatusBadRequest)
				logger.Logger().Error(
					"serverRequestBodyHash is not eq to agentRequestBodyHash",
					zap.String("agentRequestBodyHash", agentRequestBodyHash),
					zap.String("serverRequestBodyHash", serverRequestBodyHash),
					zap.Error(err),
				)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
func (m *Middleware) ResponseHashMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		if m.SecretKey != "" {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				logger.Logger().Error(
					"failed to read response",
					zap.Error(err),
				)
				return
			}
			encodedResponseBodyHash := hash.Hash(b, m.SecretKey)
			w.Header().Set(hashHeaderKey, encodedResponseBodyHash)
		}
	})
}
