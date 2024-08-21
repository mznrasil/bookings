package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mznrasil/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var appConfig *config.AppConfig
	mux := routes(appConfig)
	switch v := mux.(type) {
	case *chi.Mux:
	// pass
	default:
		t.Error(fmt.Sprintf("Type is not *chi.Mux, but %T", v))
	}
}
