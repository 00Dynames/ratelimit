// Package ratelimit provides a rate limiter for use with http requests
package ratelimit

import (
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net"
	"net/http"
	"time"
)

type rateLimiter interface {
	Allow(requests int, interval time.Duration) (bool, float64)
}

// A Limiter stores a map of *rate.limiters where each entry
// represents a different request origin, using the IP address
// as the key.
// It provides a method Limit that checks if the rate limit
// has been exceeded for the current request.
type Limiter struct {
	rate      rate.Limit
	b         int
	userLimit bool
	visitors  map[string]*rate.Limiter
}

// NewLimiter takes a limit and an integer b where
// b represents a b sized bucket of tokens that is
// re-filled at limit rps (requests per second).
// See https://en.wikipedia.org/wiki/Token_bucket for more about token buckets.
// It also takes a boolean flag userLimit that will rate limit on a
// per user basis when set to true.
// It returns an instance of a Limiter struct.
func NewLimiter(limit float64, b int, userLimit bool) *Limiter {
	return &Limiter{rate.Limit(limit), b, userLimit, make(map[string]*rate.Limiter)}
}

// Limit limits the number of requests taken by the server.
// It takes anything that implements the http.Handler interface and returns
// a handler function that checks the request rate and serves the request
// if the the limit has not been reached.
// If the limit is reached a 429 http code is returned to the requester
// with a message that indicates how long they need to wait until
// another request can be made.
// If the userLimit flag is set it fetches the limiter for  the specified user
// otherwise it fetches a limiter assigned to itself.
func (lim *Limiter) Limit(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			ip := "self"
			var err error
			if lim.userLimit {
				ip, _, err = net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					log.Println(err.Error())
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			}

			delay := lim.getVisitor(ip).Reserve().Delay()
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

// Retrieve and return the rate limiter for the current visitor if it
// already exists. Otherwise create a new rate limiter and add it to
// the visitors map, using the IP address as the key.
func (lim *Limiter) getVisitor(ip string) *rate.Limiter {
	limiter, exists := lim.visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(lim.rate, lim.b)
		lim.visitors[ip] = limiter
	}

	return limiter
}

// Rate takes a request limit and a time interval and calculates
// the rate.Limit and n inputs for rate.NewLimiter
func Rate(req int, interval time.Duration) (rate.Limit, int) {

	seconds := int(interval / time.Second)
	r := float64(req) / float64(seconds)
	return rate.Limit(r), req
}
