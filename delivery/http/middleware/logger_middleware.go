package middleware

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func LoggerMiddleware(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		h(w, r, p)
	}
}

func LoggerFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[REQUEST] %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}