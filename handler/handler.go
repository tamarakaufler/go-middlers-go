package handler

import (
	"net/http"

	"github.com/tamarakaufler/go-middlers-go/middlerware/middler3"
)

//	"github.com/tamarakaufler/go-middlers-go/middlerware/middler1"
//	"github.com/tamarakaufler/go-middlers-go/middlerware/middler2"
//	"github.com/tamarakaufler/go-middlers-go/middlerware/middler3"

type Router struct {
	Mux         http.Handler
	middlerware []middler3.Middler
}

func Apply(h http.Handler, middlers ...middler3.Middler) http.Handler {
	for _, middler := range middlers {
		h = middler(h)
	}
	return h
}
