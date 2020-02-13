package ratelimit

import (
	"sync"
	"time"
)

// Bucket struct
type Bucket struct {
	capacity int
	tokens   chan struct{}
	rate     time.Duration
	mutex    sync.Mutex
}

// NewBucket creates a new bucket instance
func NewBucket(rate time.Duration, capacity int) {

	tokens := make(chan struct{}, capacity)

	b := &Bucket{capacity, tokens, rate, sync.Mutex{}}

	go func(b *Bucket) {
		ticker := time.NewTicker(rate)
		for range ticker.C {
			b.tokens <- struct{}{}
		}
	}(b)

	return b
}

// Allow implements the Allow function from the RateLimiter interface
func (b *Bucket) Allow(requests int, interval time.Duration) (bool, float64) {
	// if you can soend a token return true, 0

	// else calculate how long until you can spend another token
	// return false, wait time
}
