package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mznrasil/bookings/internal/config"
	"github.com/mznrasil/bookings/internal/handlers"
	"github.com/mznrasil/bookings/internal/helpers"
	"github.com/mznrasil/bookings/internal/models"
	"github.com/mznrasil/bookings/internal/render"
)

const PORT = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Addr:    PORT,
		Handler: routes(&appConfig),
	}

	fmt.Println("Listening on port:", PORT)
	err = server.ListenAndServe()
	log.Fatal(err)
}

func run() error {
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

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache: ", err)
		return err
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	render.NewTemplates(&appConfig)

	repo := handlers.NewRepository(&appConfig)
	helpers.NewHelpers(&appConfig)
	handlers.NewHandlers(repo)
	return nil
}
