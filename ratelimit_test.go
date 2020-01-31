package ratelimit

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 429 response should be returned when the rate limit
// is exceeded.
func TestLimitExceeded(t *testing.T) {

	limiter := NewLimiter(0.1, 1, false)

	okHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}

	result := limiter.Limit(http.HandlerFunc(okHandler))

	var rr *httptest.ResponseRecorder

	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("GET", "https://example.com/", nil)
		rr = httptest.NewRecorder()
		result.ServeHTTP(rr, req)
	}

	assert.Equal(t, 429, rr.Code)
}

// Rate limits for requests from different remote origins
// are tracked separately
func TestLimitExceededPerUser(t *testing.T) {

	limiter := NewLimiter(0.1, 1, true)
	remote1 := "remote1:4000"
	remote2 := "remote2:4000"

	okHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}

	result := limiter.Limit(http.HandlerFunc(okHandler))

	var rr *httptest.ResponseRecorder

	// Exceed rate limit for remote1
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("GET", "https://example.com/", nil)
		req.RemoteAddr = remote1
		rr = httptest.NewRecorder()
		result.ServeHTTP(rr, req)
	}

	assert.Equal(t, 429, rr.Code)

	// Assert that rate limit for remote 2 not exceeded
	req, _ := http.NewRequest("GET", "https://example.com/", nil)
	req.RemoteAddr = remote2
	rr = httptest.NewRecorder()
	result.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}
