package middleware

import (
	"io"
	"net"
	"net/http"

	"github.com/andreevym/metric-collector/internal/logger"
	"go.uber.org/zap"
)

func (m *Middleware) TrustedSubnetMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipStr := r.Header.Get("X-Real-IP")
		if ipStr == "" {
			h.ServeHTTP(w, r)
			return
		}
		ip := net.ParseIP(ipStr)

		if !m.TrustedSubnet.Contains(ip) {
			_, err := io.WriteString(w, "Trusted Subnet Not Trusted")
			if err != nil {
				logger.Logger().Error("value can't be written", zap.Error(err))
			}
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}
