package ratelimit

import (
	"net/http"
)

type Limiter struct {
	limit int
}

func NewLimiter(limit int) *Limiter {
	return &Limiter{limit}
}

func (*Limiter) Limit(handler http.Handler) http.Handler {
	return nil
}
