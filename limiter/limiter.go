package limiter

import "time"

// Limiter define a interface para o rate limiter
type Limiter interface {
	AllowRequest(identifier string, isToken bool) bool
	BlockTimeLeft(identifier string) time.Duration
}
