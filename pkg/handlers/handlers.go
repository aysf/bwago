package handlers

import (
	"net/http"

	"github.com/aysf/bwago/pkg/models"
	"github.com/aysf/bwago/pkg/render"
)

func Home(rw http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(rw, "home.page.tmpl", &models.TemplateData{})
}

func About(rw http.ResponseWriter, r *http.Request) {

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

	// passing data to template
	render.RenderTemplate(rw, "about.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}
