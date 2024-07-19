package middleware

type Middleware struct {
	SecretKey string
	CryptoKey string
}

func NewMiddleware(secretKey string, cryptoKey string) *Middleware {
	return &Middleware{
		SecretKey: secretKey,
		CryptoKey: cryptoKey,
	}
}
