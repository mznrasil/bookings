package handlers

import (
	"net/http"

	"github.com/mznrasil/bookings/pkg/config"
	"github.com/mznrasil/bookings/pkg/models"
	"github.com/mznrasil/bookings/pkg/render"
)

var Repo *Repository

type Repository struct {
	AppConfig *config.AppConfig
}

func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{
		AppConfig: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteAddr := r.RemoteAddr
	m.AppConfig.Session.Put(r.Context(), "remote_ip", remoteAddr)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	remoteAddr := m.AppConfig.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteAddr

	//send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
