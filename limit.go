package Limit

import (
	"net/http"
)

type Limiter struct {
	rate int
}

func NewLimiter(limit int) *Limiter {
	return &Limiter{limit}
}

func (*Limiter) Limit(handler http.Handler) http.Handler {
	return nil
}
