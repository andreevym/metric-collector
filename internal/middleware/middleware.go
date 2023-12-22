package middleware

type Middleware struct {
	SecretKey string
}

func NewMiddleware(secretKey string) *Middleware {
	return &Middleware{
		SecretKey: secretKey,
	}
}
