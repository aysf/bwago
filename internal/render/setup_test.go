package render

import (
	"encoding/gob"
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
	testApp.InProduction = true

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
