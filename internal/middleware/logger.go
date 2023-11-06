package middleware

import (
	"io"
	"net/http"
	"time"

	"github.com/andreevym/metric-collector/internal/logger"
	"go.uber.org/zap"
)

func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		h.ServeHTTP(w, r)
		end := time.Now()

		if r == nil {
			return
		}

		logger.Log.Info(
			"request",
			zap.String("method", r.Method),
			zap.String("URI", r.RequestURI),
			zap.Duration("duration", end.Sub(start)),
		)

		if r.Response == nil {
			return
		}
		defer func() {
			err := r.Body.Close()
			if err != nil {
				logger.Log.Error(err.Error())
			}
		}()
		bytes, _ := io.ReadAll(r.Body)

		logger.Log.Info(
			"response",
			zap.Int("status", r.Response.StatusCode),
			zap.Int("status", len(bytes)),
		)
	})
}
