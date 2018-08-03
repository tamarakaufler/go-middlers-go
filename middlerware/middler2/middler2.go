package middler2

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Middler func(http.Handler) http.Handler

// Apply runs a sequence of middlers starting from the last one provided.
// An example usage:
// 			http.Handle("/hi", Wrap(hiHandler, LoggingMiddler(logger), TracingMiddler()))
func Apply(h http.Handler, middlers ...Middler) http.Handler {
	for _, middler := range middlers {
		h = middler(h)
	}
	return h
}

// ReverseApply applies a sequence of middlers starting from the first one provided.
// An example usage:
// 			http.Handle("/hi", ReverseWrap(hiHandler, LoggingMiddler(logger), TracingMiddler()))
func ReverseApply(h http.Handler, middlers ...Middler) http.Handler {
	l := len(middlers)
	for i := l - 1; i >= 0; i-- {
		h = middlers[i](h)
	}
	return h
}

// LoggingMiddler is logging middleware. An example usage:
// http.Handle("/", LoggingMiddler(logger)(routeHandler))
func LoggingMiddler(logger *log.Logger) Middler {
	logger.Println("Initialised LoggingMiddler2")

	// This is the returned Middler --------^
	return func(h http.Handler) http.Handler {

		// This is the returned http.Handler ^
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("Logging in LoggingMiddler before")
			logger.Printf("\t%s", r.RequestURI)

			defer logger.Println("Logging in LoggingMiddler after")
			h.ServeHTTP(w, r)
		})
	}
}

func TracingMiddler() Middler {
	log.Println("Initialised TracingMiddler2")

	return func(h http.Handler) http.Handler {

		requestID := rand.Int63()

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Println("Logging in TracingMiddler before")
			log.Printf("\t%s", r.RequestURI)

			defer log.Println("Logging in TracingMiddler after")

			log.Printf("TracingMiddler: requestID %d: start %s", requestID, start)
			defer log.Printf("TracingMiddler: requestID %d: took %s", requestID, time.Since(start))

			h.ServeHTTP(w, r)
		})
	}
}
