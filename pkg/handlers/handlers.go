package handlers

import (
	"net/http"

	"github.com/aysf/bwago/pkg/config"
	"github.com/aysf/bwago/pkg/models"
	"github.com/aysf/bwago/pkg/render"
)

// Repo the repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home handles home page
func (m *Repository) Home(rw http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(rw, "home.page.tmpl", &models.TemplateData{})
}

// About handles about page
func (m *Repository) About(rw http.ResponseWriter, r *http.Request) {

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	// peform some logic
	// get data from struct
	data := make(map[string]interface{})
	data["satu"] = struct {
		Nama string
		Umur int
	}{
		"Mina",
		4,
	}

	// get data from string map
	stringMap := make(map[string]string)
	stringMap["hobby"] = "drawing"

	stringMap["remote_ip"] = remoteIP

	// passing data to template
	render.RenderTemplate(rw, "about.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// About handles about page
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, "contact.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and display form
func (m *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, "make-reservation.page.tmpl", &models.TemplateData{})
}

// Generals renders the room page
func (m *Repository) Generals(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, "majors.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, "search-availability.page.tmpl", &models.TemplateData{})
}
