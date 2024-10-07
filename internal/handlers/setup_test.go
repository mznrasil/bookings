package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/mznrasil/bookings/internal/config"
	"github.com/mznrasil/bookings/internal/models"
	"github.com/mznrasil/bookings/internal/render"
)

var appConfig config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates/"

func getRoutes() http.Handler {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	appConfig.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache: ", err)
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = true

	repo := NewRepository(&appConfig)
	NewHandlers(repo)
	render.Template(&appConfig)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", http.HandlerFunc(Repo.Home))
	mux.Get("/about", http.HandlerFunc(Repo.About))
	mux.Get("/contact", http.HandlerFunc(Repo.Contact))
	mux.Get("/generals-quarters", http.HandlerFunc(Repo.Generals))
	mux.Get("/majors-suite", http.HandlerFunc(Repo.Majors))

	mux.Get("/make-reservation", http.HandlerFunc(Repo.MakeReservation))
	mux.Post("/make-reservation", http.HandlerFunc(Repo.PostMakeReservation))
	mux.Get("/reservation-summary", http.HandlerFunc(Repo.ReservationSummary))

	mux.Get("/search-availability", http.HandlerFunc(Repo.SearchAvailability))
	mux.Post("/search-availability", http.HandlerFunc(Repo.PostSearchAvailability))
	mux.Post("/search-availability-json", http.HandlerFunc(Repo.SearchAvailabilityJSON))

	mux.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	return mux
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
