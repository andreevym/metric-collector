package middleware

import (
	"io"
	"net"
	"net/http"
)

func (m *Middleware) TrustedSubnetMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipStr := r.Header.Get("X-Real-IP")
		ip := net.ParseIP(ipStr)

		if !m.TrustedSubnet.Contains(ip) {
			_, err := io.WriteString(w, "Trusted Subnet Not Trusted")
			if err != nil {
				panic(err)
			}
			w.WriteHeader(http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}
