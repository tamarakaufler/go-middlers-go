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

	piggyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		time.Sleep(2 * time.Second)
		log.Printf("%s", r.RequestURI)
		w.Write([]byte("Had a peep in my piggybank :)"))
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

	auth := &Auth{
		Username: "Tamara",
	}
	authMiddling := &middler3.Middling{}
	authMiddling.Use(middler3.LoggingMiddler(logger))
	authMiddling.Use(auth.Middlerware())

	// routing
	// -------

	// applying just one middler
	http.Handle("/", rootMiddling.Apply(rootHandler))

	// applying a sequence of middlers
	http.Handle("/hi", hiMiddling.Apply(hiHandler))

	// applying just one middler, different than for the root route
	http.Handle("/bye", byeMiddling.Apply(byeHandler))

	// applying one Middler and one custom middler
	http.Handle("/piggybank", authMiddling.Apply(piggyHandler))

	// Some browsers request /favicon.ico, even though the server does not.
	// The following is added to avoid the annoyance through serving an empty
	// favicon to avoid the middleware effect
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})

	// spin up the server
	if err := http.ListenAndServe(":7777", nil); err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %q: %s\n", 7777, err)
	}
}
