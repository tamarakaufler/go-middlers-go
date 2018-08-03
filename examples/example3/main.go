package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tamarakaufler/go-middlers-go/middlerware/middler3"
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
	logger := log.New(os.Stdout, "with-middler3: ", log.LstdFlags)

	rootMiddling := &middler3.Middling{}
	rootMiddling.Use(middler3.LoggingMiddler(logger))

	hiMiddling := &middler3.Middling{}
	hiMiddling.Use(middler3.LoggingMiddler(logger))
	hiMiddling.Use(middler3.TracingMiddler())

	byeMiddling := &middler3.Middling{}
	byeMiddling.Use(middler3.TracingMiddler())

	// routing
	// -------

	// applying just one middler
	http.Handle("/", rootMiddling.Apply(rootHandler))

	// applying a sequence of middlers
	http.Handle("/hi", hiMiddling.Apply(hiHandler))

	// applying just one middler, different than for the root route
	http.Handle("/bye", byeMiddling.Apply(byeHandler))

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
