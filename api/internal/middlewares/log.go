package middlewares

import (
	"log"
	"net/http"
	"time"
)

// Log proxies requests and logging their execution information
func Log(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		executionTime := time.Since(startTime)
		log.Printf("%s %s %s", r.Method, r.URL.Path, executionTime.String())
		next.ServeHTTP(w, r)
	}
}
