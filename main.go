package main

import (
	"time"
	//	"github.com/tamarakaufler/go-middlers-go/middlerware/middler1"
	//	"github.com/tamarakaufler/go-middlers-go/middlerware/middler2"
	//	"github.com/tamarakaufler/go-middlers-go/middlerware/middler3"

	"net/http"

	"github.com/tamarakaufler/go-middlers-go/middlerware/middler1"
)

func main() {

	m1 := middler1.Middlers{middler1.LoggingMiddler, middler1.TracingMiddler}
	http.HandleFunc("/", m1.Apply(rootHandler))
	http.HandleFunc("/hi", hiHandler)
	http.HandleFunc("/bye", byeHandler)

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	time.Sleep(1 * time.Second)
	w.Write([]byte("I am the ROOT of the matter"))
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
