package middleware

import (
	"io"
	"net/http"
	"time"

	"github.com/andreevym/metric-collector/internal/logger"
	"go.uber.org/zap"
)

// RequestLoggerMiddleware returns an HTTP middleware that logs information about incoming requests
// and outgoing responses. It records the HTTP method, URI, request duration, response status code,
// and response body length.
//
// Parameters:
//   - h: The HTTP handler to be wrapped by the middleware.
//
// Returns:
//   - http.Handler: An HTTP handler that logs request and response details.
//
// Example:
//
//	// Create a new middleware instance
//	middleware := NewMiddleware("my-secret-key")
//
//	// Wrap an existing HTTP handler with the RequestLoggerMiddleware
//	wrappedHandler := middleware.RequestLoggerMiddleware(myHandler)
func (m *Middleware) RequestLoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time of the request
		start := time.Now()

		// Call the next HTTP handler in the chain
		h.ServeHTTP(w, r)

		// Record the end time of the request
		end := time.Now()

		if r == nil {
			return
		}

		// Log information about the incoming request
		logger.Logger().Info(
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
				logger.Logger().Error(err.Error())
			}
		}()
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Logger().Error("read body error", zap.Error(err))
			return
		}

		// Log information about the outgoing response
		logger.Logger().Info(
			"response",
			zap.Int("status", r.Response.StatusCode),
			zap.Int("status", len(bytes)),
		)
	})
}
