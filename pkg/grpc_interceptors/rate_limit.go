package grpc_interceptors

import (
	"context"
	"errors"

	"golang.org/x/time/rate"
)

var errRateLimitExceeede = errors.New("rate limit exceeded")

type rateLimiter struct {
	limiter *rate.Limiter
}

// NewRateLimiter -
func NewRateLimiter(r rate.Limit, b int) *rateLimiter {
	return &rateLimiter{
		limiter: rate.NewLimiter(r, b),
	}
}

func (rl *rateLimiter) Limit(ctx context.Context) error {
	if !rl.limiter.Allow() {
		return errRateLimitExceeede
	}

	return nil
}
