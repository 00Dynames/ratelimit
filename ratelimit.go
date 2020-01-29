// Package ratelimit provides a rate limiter for use with http requests
package ratelimit

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
)

// A Limiter wraps a Limiter struct from the rate package
// and provides a method Limit that checks if the rate limit
// has been exceeded for the current request.
type Limiter struct {
	l *rate.Limiter
}

// NewLimiter takes a limit that represents requests per second (rps)
// and returns an instance of a Limiter struct.
func NewLimiter(limit int) *Limiter {
	return &Limiter{rate.NewLimiter(rate.Limit(limit), 4)}
}

// Limit takes anything that implements the http.Handler interface and returns
// a handler function that serves the request.
func (limiter *Limiter) Limit(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			delay := limiter.l.Reserve().Delay()
			if delay != 0 {
				http.Error(
					w,
					fmt.Sprintf("Rate limit exceeded. Try again in %v", delay),
					http.StatusTooManyRequests,
				)
				return
			}

			handler.ServeHTTP(w, r)
		},
	)
}
