package api

import "crypto/rsa"

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Limiter struct {
		RPS     float64
		Burst   int
		Enabled bool
	}
	Cors struct {
		TrustedOrigins []string
	}
	Auth struct {
		SigningKey               *rsa.PublicKey
		TokenExpirationInMinutes int
	}
	LlmApiKey string
}
