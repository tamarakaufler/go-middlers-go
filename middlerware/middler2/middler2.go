package middler2

import (
	"log"
	"net/http"
)

type Middler func(http.Handler) http.Handler

// LoggingMiddler is logging middleware. An example usage:
// http.Handle("/", LoggingMiddler(logger)(routeHandler))
func LoggingMiddler(logger log.Logger) Middler {
	logger.Println("Initialised LoggingMiddler2")

	// This is the returned Middler --------^
	return func(h http.Handler) http.Handler {

		// This is the returned http.Handler ^
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("Logging before")
			defer logger.Println("Logging after")
			h.ServeHTTP(w, r)
		})
	}
}
