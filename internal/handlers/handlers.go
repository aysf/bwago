package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aysf/bwago/internal/config"
	"github.com/aysf/bwago/internal/forms"
	"github.com/aysf/bwago/internal/models"
	"github.com/aysf/bwago/internal/render"
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

	render.RenderTemplate(rw, r, "home.page.tmpl", &models.TemplateData{})
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
	render.RenderTemplate(rw, r, "about.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// About handles about page
func (m *Repository) Contact(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "contact.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and display form
func (m *Repository) Reservation(rw http.ResponseWriter, r *http.Request) {
	emptyReservation := models.Reservation{}
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	forms := forms.New(r.PostForm)

	// forms.Has("first_name", r)

	forms.Required("first_name", "last_name", "email")
	forms.MinLength("first_name", 3, r)
	forms.IsEmail("email")

	if !forms.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(rw, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: forms,
			Data: data,
		})
		return

	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(rw, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders search availability form
func (m *Repository) Availability(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability accept post request from search availability form
func (m *Repository) PostAvailability(rw http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	rw.Write([]byte(fmt.Sprintf("the start date is %s, and the end date is %s", start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJson handles request for availability and send JSON response
func (m *Repository) AvailabilityJson(rw http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      true,
		Message: "Available !",
	}

	out, err := json.MarshalIndent(resp, "", "	")
	if err != nil {
		log.Println(err)
	}

	log.Println(string(out))

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(out)

}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cannot get session from template")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		log.Println("cannot get item from session")
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})

}
