package ratelimit

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
)

type Limiter struct {
	l *rate.Limiter
}

func NewLimiter(limit int) *Limiter {
	fmt.Println(limit)
	return &Limiter{rate.NewLimiter(rate.Limit(limit), 4)}
}

func (limiter *Limiter) Limit(handler http.Handler) http.Handler {
	limiter.l.Limit()
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
