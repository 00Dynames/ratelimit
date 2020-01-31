# ratelimit

Package ratelimit provides a rate limiter for use with http requests.
Further documentation can be found at https://godoc.org/github.com/00Dynames/ratelimit.

Usage example
-------------

```golang
package main

import (
  "github.com/00Dynames/ratelimit"
  "net/http"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", okHandler)

  // limiter restricts requests to 1 request per second
  limiter := ratelimit.NewLimiter(1, 1, false)

  http.ListenAndServe(":4000", limiter.Limit(mux))
}

func okHandler(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("OK"))
}
```
