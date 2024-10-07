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
	"github.com/mznrasil/bookings/internal/driver"
	"github.com/mznrasil/bookings/internal/handlers"
	"github.com/mznrasil/bookings/internal/helpers"
	"github.com/mznrasil/bookings/internal/models"
	"github.com/mznrasil/bookings/internal/render"
)

const PORT = ":8080"

var (
	appConfig config.AppConfig
	session   *scs.SessionManager
)

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	server := http.Server{
		Addr:    PORT,
		Handler: routes(&appConfig),
	}

	fmt.Println("Listening on port:", PORT)
	err = server.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.Restriction{})

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

	// Connect to a database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=rasil password=")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	fmt.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache: ", err)
		return nil, err
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	render.NewRenderer(&appConfig)

	repo := handlers.NewRepository(&appConfig, db)
	helpers.NewHelpers(&appConfig)
	handlers.NewHandlers(repo)
	return db, nil
}
