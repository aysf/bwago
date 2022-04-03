package main

import (
	"log"
	"net/http"

	"github.com/justinas/nosurf"
)

// ConsoleLog writes log in the console
func ConsoleLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Println("hit the page")
		next.ServeHTTP(rw, r)
	})
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and save in every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
