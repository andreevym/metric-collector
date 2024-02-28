package middleware

import (
	"io"
	"net/http"

	"github.com/andreevym/metric-collector/internal/hash"
	"github.com/andreevym/metric-collector/internal/logger"
	"go.uber.org/zap"
)

const HashHeaderKey = "HashSHA256"

// RequestHashMiddleware returns an HTTP middleware that verifies the hash of the request body
// sent by the client against the calculated hash on the server side, if the client provides a hash.
// If the provided hash does not match the calculated hash, the middleware returns a Bad Request status code.
// The calculated hash is based on the secret key provided in the middleware.
// If no hash is provided by the client or if the hash verification fails, the middleware proceeds to the next handler.
//
// Parameters:
//   - h: The HTTP handler to be wrapped by the middleware.
//
// Returns:
//   - http.Handler: An HTTP handler that verifies the hash of the request body against the calculated hash.
//
// Example:
//
//	// Create a new middleware instance
//	middleware := NewMiddleware("my-secret-key")
//
//	// Wrap an existing HTTP handler with the RequestHashMiddleware
//	wrappedHandler := middleware.RequestHashMiddleware(myHandler)
func (m *Middleware) RequestHashMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the hash of the request body sent by the client from the request headers
		agentRequestBodyHash := r.Header.Get(HashHeaderKey)
		if agentRequestBodyHash != "" {
			// Read the entire request body
			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				// If an error occurs while reading the request body, set a Bad Request status code
				// and log the error
				w.WriteHeader(http.StatusBadRequest)
				logger.Logger().Error(
					"could not read request body when agent hash sha 256 is defined",
					zap.String("agentRequestBodyHash", agentRequestBodyHash),
					zap.Error(err),
				)
				return
			}

			// Calculate the hash of the request body using the specified secret key
			serverRequestBodyHash := hash.EncodeHash(bytes, m.SecretKey)

			// Compare the hash provided by the client with the calculated hash on the server side
			if serverRequestBodyHash != agentRequestBodyHash {
				// If the hashes do not match, set a Bad Request status code
				// and log the error
				w.WriteHeader(http.StatusBadRequest)
				logger.Logger().Error(
					"serverRequestBodyHash is not equal to agentRequestBodyHash",
					zap.String("agentRequestBodyHash", agentRequestBodyHash),
					zap.String("serverRequestBodyHash", serverRequestBodyHash),
				)
				return
			}
		}

		// Proceed to the next HTTP handler in the chain
		h.ServeHTTP(w, r)
	})
}

// ResponseHashMiddleware returns an HTTP middleware that calculates a hash of the response body
// and includes it in the response headers if a secret key is provided.
// The calculated hash is encoded using the specified secret key.
// If the secret key is empty, the middleware does not modify the response.
//
// Parameters:
//   - h: The HTTP handler to be wrapped by the middleware.
//
// Returns:
//   - http.Handler: An HTTP handler that includes the calculated hash in the response headers.
//
// Example:
//
//	// Create a new middleware instance
//	middleware := NewMiddleware("my-secret-key")
//
//	// Wrap an existing HTTP handler with the ResponseHashMiddleware
//	wrappedHandler := middleware.ResponseHashMiddleware(myHandler)
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
			encodedResponseBodyHash := hash.EncodeHash(b, m.SecretKey)
			w.Header().Set(HashHeaderKey, encodedResponseBodyHash)
		}
	})
}
