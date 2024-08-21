package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mznrasil/bookings/internal/config"
	"github.com/mznrasil/bookings/internal/handlers"
)

func routes(appConfig *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/contact", http.HandlerFunc(handlers.Repo.Contact))
	mux.Get("/generals-quarters", http.HandlerFunc(handlers.Repo.Generals))
	mux.Get("/majors-suite", http.HandlerFunc(handlers.Repo.Majors))

	mux.Get("/make-reservation", http.HandlerFunc(handlers.Repo.MakeReservation))
	mux.Post("/make-reservation", http.HandlerFunc(handlers.Repo.PostMakeReservation))
	mux.Get("/reservation-summary", http.HandlerFunc(handlers.Repo.ReservationSummary))

	mux.Get("/search-availability", http.HandlerFunc(handlers.Repo.SearchAvailability))
	mux.Post("/search-availability", http.HandlerFunc(handlers.Repo.PostSearchAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(handlers.Repo.SearchAvailabilityJSON))

	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	return mux
}
