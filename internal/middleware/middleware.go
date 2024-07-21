package middleware

import "net"

type Middleware struct {
	SecretKey     string
	CryptoKey     string
	TrustedSubnet *net.IPNet
}

func NewMiddleware(secretKey string, cryptoKey string, trustedSubnet *net.IPNet) *Middleware {
	return &Middleware{
		SecretKey:     secretKey,
		CryptoKey:     cryptoKey,
		TrustedSubnet: trustedSubnet,
	}
}
