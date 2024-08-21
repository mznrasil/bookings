package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mznrasil/bookings/internal/config"
	"github.com/mznrasil/bookings/internal/models"
)

var session *scs.SessionManager
var testAppConfig config.AppConfig

func TestMain(m *testing.M) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})

	// change this to true when in production
	testAppConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testAppConfig.Session = session

	appConfig = &testAppConfig

	os.Exit(m.Run())
}

type myWriter struct{}

func (mw myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (mw myWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (mw myWriter) WriteHeader(statusCode int) {}
