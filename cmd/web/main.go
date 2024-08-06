package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mznrasil/bookings/pkg/config"
	"github.com/mznrasil/bookings/pkg/handlers"
	"github.com/mznrasil/bookings/pkg/render"
)

const PORT = ":8080"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {
	// change this to true when in production
	appConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache: ", err)
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	render.NewTemplates(&appConfig)

	repo := handlers.NewRepository(&appConfig)
	handlers.NewHandlers(repo)

	server := http.Server{
		Addr:    PORT,
		Handler: routes(&appConfig),
	}

	fmt.Println("Listening on port:", PORT)
	err = server.ListenAndServe()
	log.Fatal(err)
}
