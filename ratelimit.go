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

// NewLimiter takes a limit and an integer b where
// b represents a b sized bucket of tokens that is
// re-filled at limit rps (requests per second).
// See https://en.wikipedia.org/wiki/Token_bucket for more about token buckets.
// It returns an instance of a Limiter struct.
func NewLimiter(limit float64, b int) *Limiter {
	return &Limiter{rate.NewLimiter(rate.Limit(limit), b)}
}

// Limit limits the number of requests taken by the server.
// It takes anything that implements the http.Handler interface and returns
// a handler function that checks the request rate and serves the request
// if the the limit has not been reached.
// If the limit is reached a 429 http code is returned to the requester
// with a message that indicates how long they need to wait until
// another request can be made.
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
