package main

import (
	"log"
	"net/http"
	"time"

	"github.com/tamarakaufler/go-middlers-go/middlerware/middler1"
)

func main() {

	router := http.NewServeMux()

	m1 := middler1.Middlers{middler1.TracingMiddler, middler1.LoggingMiddler}
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/hi", hiHandler)
	router.HandleFunc("/bye", byeHandler)

	server := &http.Server{
		Addr:    ":7777",
		Handler: m1.Apply(router),
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %q: %s\n", 7777, err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	time.Sleep(1 * time.Second)

	msg := "I am the ROOT of the matter"
	log.Print(msg)
	w.Write([]byte(msg))
}

func hiHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	time.Sleep(2 * time.Second)
	w.Write([]byte("Hi world. You are challenging."))
}

func byeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	time.Sleep(3 * time.Second)
	w.Write([]byte("Bye world. I am going to sleep"))
}
