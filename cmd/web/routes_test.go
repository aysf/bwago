package main

import (
	"testing"

	"github.com/aysf/bwago/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {

	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Errorf("route output is not http.Handler, but %T", v)
	}
}
