package handlers

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/aysf/bwago/internal/config"
	"github.com/aysf/bwago/internal/models"
	"github.com/aysf/bwago/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var templateDir = "./../../templates"

// var functions = template.FuncMap{}

func TestMain(m *testing.M) {
	// copy from main func
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	app.InProduction = true

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infolog

	errorlog := log.New(os.Stdout, "ERRO\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorlog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	// app.UseCache = true

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)

	listenForMail()

	tc, err := createTestTemplateCache()
	if err != nil {
		log.Panic("error loading template cache ", err)
	}

	app.TemplateCache = tc
	app.UseCache = app.InProduction

	Repo := NewTestRepo(&app)
	NewHandlers(Repo)
	render.NewRenderer(&app)

	// before exit the func, run the test
	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

func getRoute() http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(consoleLog)
	mux.Use(sessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJson)

	mux.Get("/choose-room/{id}", Repo.ChooseRoom)
	mux.Get("/book-room", Repo.BookRoom)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fs))

	return mux

}

// consoleLog writes log in the console
func consoleLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Println("hit the test page")
		next.ServeHTTP(rw, r)
	})
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// sessionLoad loads and save in every request
func sessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// createTestTemplateCache creates template cache in a map
func createTestTemplateCache() (map[string]*template.Template, error) {

	// tc is template cache
	tc := map[string]*template.Template{}

	// getting all page templates
	pages, err := filepath.Glob(templateDir + "/*.page.gohtml")
	if err != nil {
		return tc, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		// ts is template set
		ts, err := template.ParseFiles(page)
		if err != nil {
			return tc, err
		}

		// getting all layout templates
		layouts, err := filepath.Glob(templateDir + "/*.layout.gohtml")
		if err != nil {
			return tc, err
		}

		if len(layouts) > 0 {

			ts, err = ts.ParseGlob(templateDir + "/*.layout.gohtml")
			if err != nil {
				return tc, err
			}
		}

		tc[name] = ts
	}

	return tc, nil
}
