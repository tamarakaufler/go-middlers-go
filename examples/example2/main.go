package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tamarakaufler/go-middlers-go/middlerware/middler2"
)

func main() {

	// Handler definitions
	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		time.Sleep(1 * time.Second)

		msg := "I am the ROOT of the matter"
		log.Print(msg)
		w.Write([]byte(msg))
	})

	hiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		time.Sleep(2 * time.Second)
		w.Write([]byte("Hi world. You are challenging."))
	})

	byeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		time.Sleep(3 * time.Second)
		w.Write([]byte("Bye world. I am going to sleep"))
	})

	// for middleware
	logger := log.New(os.Stdout, "with-middler2: ", log.LstdFlags)

	// routing

	// applying just one middler
	http.Handle("/", middler2.LoggingMiddler(logger)(rootHandler))

	// applying a sequence of middlers - they are applied in the opposite sequence
	// to how they were provided
	http.Handle("/hi", middler2.Wrap(hiHandler, middler2.LoggingMiddler(logger), middler2.TracingMiddler()))

	// applying a sequence of middlers - they are applied in the same sequence
	// as they were provided
	http.Handle("/bye", middler2.ReverseWrap(byeHandler, middler2.LoggingMiddler(logger), middler2.TracingMiddler()))

	// some browsers request /favicon.ico, even though the server is not. This is to avoid
	// the annoyance through serving an empty favicon to avoid the middleware effect
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})

	// spin up the server
	if err := http.ListenAndServe(":7777", nil); err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %q: %s\n", 7777, err)
	}
}
