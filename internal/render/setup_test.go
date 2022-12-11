package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/aysf/bwago/internal/config"
	"github.com/aysf/bwago/internal/models"
)

var (
	session *scs.SessionManager
	testApp config.AppConfig
)

func TestMain(m *testing.M) {

	// copy from main func
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testApp.InProduction = false

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infolog

	errorlog := log.New(os.Stdout, "ERRO\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorlog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session
	// app.UseCache = true
	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (m *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (m *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

func (m *myWriter) WriteHeader(statusCode int) {

}
