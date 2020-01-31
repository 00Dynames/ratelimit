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

	limiter := NewLimiter(0.1, 1)

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
