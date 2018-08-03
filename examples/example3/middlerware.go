package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/tamarakaufler/go-middlers-go/middlerware/middler3"
)

type Auth struct {
	Username string
}

func (au *Auth) Middlerware() middler3.Middler {
	log.Println("\t>> My Auth Middleware")

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("No Authorization header")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			authParts := strings.Split(authHeader, " ")
			if authParts[0] != "Bearer" {
				log.Println("Unauthorised - incorrect Authorization header")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// get the stored token for the username from somewhere
			stored := "12345"
			if authParts[1] != stored {
				log.Println("Unauthorised - Forbidden")
				w.WriteHeader(http.StatusForbidden)
				return
			}

			log.Println("Authorised - before")
			defer log.Println("Authorised - after")

			h.ServeHTTP(w, r)
		})
	}
}
