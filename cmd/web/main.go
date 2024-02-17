package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/aysf/bwago/internal/config"
	"github.com/aysf/bwago/internal/driver"
	"github.com/aysf/bwago/internal/handlers"
	"github.com/aysf/bwago/internal/helpers"
	"github.com/aysf/bwago/internal/models"
	"github.com/aysf/bwago/internal/render"
)

const portNumber = ":9000"

var app config.AppConfig
var session *scs.SessionManager
var infolog *log.Logger
var errorlog *log.Logger

// main is the main application funcion
func main() {

	db, err := Run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("starting mail listener...")
	listenForMail()

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

func Run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")

	dbURL := flag.String("dburl", "", "Database URL")
	// dbSSL := flag.String("dbssl", "disable", "Database SSL settings (disable, prefer, require)")

	flag.Parse()

	// change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache

	infolog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infolog

	errorlog = log.New(os.Stdout, "ERRO\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorlog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("connecting to database...")
	log.Println("getting db url from environment...")
	db_url := os.Getenv("DBURL")
	if db_url == "" {
		log.Println("missing required DBURL ENV")
	}

	if *dbURL != "" {
		log.Println("setting database url from flag")
		db_url = *dbURL
	}
	log.Println("database url:", db_url)

	db, err := driver.ConnectSQL(db_url)
	if err != nil {
		log.Fatal("cannot connect to database! Dying...")
	}

	if err := db.SQL.Ping(); err != nil {
		log.Println("connected to database!")
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Panic("error loading template cache ", err)
		return nil, err
	}

	app.TemplateCache = tc

	routes(&app)

	Repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(Repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
