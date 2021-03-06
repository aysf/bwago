package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/aysf/bwago/internal/config"
	"github.com/aysf/bwago/internal/handlers"
	"github.com/aysf/bwago/internal/models"
	"github.com/aysf/bwago/internal/render"
)

const portNumber = ":9000"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application funcion
func main() {

	err := Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("application running on port %s\n", portNumber)
	// http.ListenAndServe(":8080", nil)

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Panic("error loading template cache ", err)
		return err
	}

	app.TemplateCache = tc
	app.UseCache = app.InProduction

	Repo := handlers.NewRepo(&app)
	handlers.NewHandlers(Repo)

	render.NewTemplates(&app)

	routes(&app)

	// http.HandleFunc("/", handlers.Home)
	// http.HandleFunc("/about", handlers.About)

	// http.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	return nil
}
