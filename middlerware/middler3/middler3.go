package middler3

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Middler func(http.Handler) http.Handler

func Apply(h http.Handler, middlers ...Middler) http.Handler {

	for _, middler := range middlers {
		h = middler(h)
	}
	return h
}

func LoggingMiddler(logger log.Logger) Middler {
	logger.Println("Initialised LoggingMiddler3")

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("Logging before")
			defer logger.Println("Logging after")
			h.ServeHTTP(w, r)
		})
	}
}

func TracingMiddler() Middler {
	log.Println("Initialised TracingMiddler3")

	requestID := rand.Int63()

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Printf("requestID %d: start %s", requestID, start)
			defer log.Printf("requestID %d: took %s", requestID, time.Since(start))

			h.ServeHTTP(w, r)
		})
	}
}
