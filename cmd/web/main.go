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

	// -> sample using std.lib for sending email
	// from := "ananto@here.com"
	// auth := smtp.PlainAuth("", from, "", "localhost")
	// err = smtp.SendMail("localhost:1025", auth, from, []string{"you@there.com"}, []byte("hello world :)"))
	// if err != nil {
	// 	log.Println(err)
	// }

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
	dbHost := flag.String("dbhost", "localhost", "Database Host")
	dbName := flag.String("dbname", "", "Database Name")
	dbUser := flag.String("dbuser", "", "Database User")
	dbPass := flag.String("dbpass", "", "Database Password")
	dbPort := flag.String("dbport", "5432", "Database Port")
	dbSSL := flag.String("dbssl", "disable", "Database SSL settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("missing required flags")
		os.Exit(1)
	}

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
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("cannot connect to database! Dying...")
	}
	log.Println("connected to database!")

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

	// http.HandleFunc("/", handlers.Home)
	// http.HandleFunc("/about", handlers.About)
	// http.Handle("/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	return db, nil
}
