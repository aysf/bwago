package main

import (
	"net/http"

	"github.com/aysf/bwago/pkg/config"
	"github.com/aysf/bwago/pkg/handlers"
	"github.com/go-chi/chi/v5"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Get("/", handlers.Home)
	mux.Get("/about", handlers.About)

	return mux
}
