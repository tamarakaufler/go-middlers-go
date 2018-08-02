package middler1

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Middler func(http.Handler) http.Handler
type Middlers []Middler

// Apply applies middleware from the last one provided, ie
// the one closest to the route handler
func (middlers Middlers) Apply(h http.Handler) http.Handler {
	if len(middlers) == 0 {
		return h
	}

	last := len(middlers) - 1
	return middlers[:last].Apply(middlers[last](h))
}

func LoggingMiddler(h http.Handler) http.Handler {
	log.Println("Initialised LoggingMiddler1")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Logging before")
		defer log.Println("Logging after")
		h.ServeHTTP(w, r)
	})
}

func TracingMiddler(h http.Handler) http.Handler {
	log.Println("Initialised TracingMiddler1")

	requestID := rand.Int63()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("requestID %d: start %s", requestID, start)
		defer log.Printf("requestID %d: took %s", requestID, time.Since(start))

		h.ServeHTTP(w, r)
	})
}
