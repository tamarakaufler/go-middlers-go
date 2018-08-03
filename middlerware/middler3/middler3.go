package middler3

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Middler is a function which receives an http.Handler and returns another http.Handler.
// Typically, the returned handler is a closure which does something with the http.ResponseWriter and http.Request passed to it, and then calls the handler passed as parameter to the MiddlewareFunc.
// (comment from gorilla middleware)
type Middler func(http.Handler) http.Handler

// User creates an empty Middling instance, which is then populated with Middlers
// using the Use method
type Middling struct {
	middlers []middlerware
}

// whatever implements the middlerware interface qualifies as middleware
type middlerware interface {
	Middlerware(handler http.Handler) http.Handler
}

// Middlerware allows Middler to implement the middleware interface
// and so, to be of type middleware.
// Other middleware implementations (not of the Middler type) must
// implement their own Middleware method.
func (m Middler) Middlerware(handler http.Handler) http.Handler {
	log.Println("\t>> Middleware")
	return m(handler)
}

// Use appends a Middler to the Middlers chain. Middlers can intercept and/or
// modify requests and/or responses. They are executed in the order that they
// are applied. Different Middlings with Middler chains can be created to be applied to
// different routes.
func (middling *Middling) Use(m middlerware) {
	middling.middlers = append(middling.middlers, m)
	log.Printf(">> Use: %d", len(middling.middlers))
}

// Apply runs middlers in the order they were provided
func (middling *Middling) Apply(h http.Handler) http.Handler {
	log.Printf(">> Apply: %d", len(middling.middlers))

	if len(middling.middlers) == 0 {
		return h
	}

	for i := len(middling.middlers) - 1; i >= 0; i-- {
		log.Printf("\t>> Apply: i=%d", i)
		h = middling.middlers[i].Middlerware(h)
	}
	return h
}

// Middlers from the box -------------------------------------------------

func LoggingMiddler(logger *log.Logger) Middler {
	logger.Println("Initialised LoggingMiddler3")

	return func(h http.Handler) http.Handler {

		// This handler is run when a route is requested
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("Logging in LoggingMiddler before")
			h.ServeHTTP(w, r)
			logger.Println("Logging in LoggingMiddler after")
		})
	}
}

func TracingMiddler() Middler {
	log.Println("Initialised TracingMiddler3")

	requestID := rand.Int63()

	return func(h http.Handler) http.Handler {

		// This handler is run when a route is requested
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Println("Logging in TracingMiddler before")
			defer log.Println("Logging in TracingMiddler after")

			log.Printf("requestID %d: start %s", requestID, start)
			h.ServeHTTP(w, r)
			log.Printf("requestID %d: took %s", requestID, time.Since(start))
		})
	}
}
