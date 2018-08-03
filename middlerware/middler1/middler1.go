package middler1

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Middler func(http.Handler) http.Handler
type Middlers []Middler

// Apply applies middleware from the first one provided, ie
// sequentially, ie the middleware is run from the last one provided
// to the first one.
func (middlers Middlers) Apply(h http.Handler) http.Handler {
	if len(middlers) == 0 {
		return h
	}
	return middlers[1:].Apply(middlers[0](h))
}

func LoggingMiddler(h http.Handler) http.Handler {
	log.Println("Initialised LoggingMiddler1")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Logging in LoggingMiddler before")
		log.Printf("%s", r.RequestURI)

		defer log.Println("Logging in LoggingMiddler after")
		h.ServeHTTP(w, r)
	})
}

func TracingMiddler(h http.Handler) http.Handler {
	log.Println("Initialised TracingMiddler1")

	requestID := rand.Int63()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Println("Logging in TracingMiddler before")
		log.Printf("%s", r.RequestURI)

		defer log.Println("Logging in TracingMiddler after")

		log.Printf("TracingMiddler: requestID %d: start %s", requestID, start)
		defer log.Printf("TracingMiddler: requestID %d: took %s", requestID, time.Since(start))

		h.ServeHTTP(w, r)
	})
}
