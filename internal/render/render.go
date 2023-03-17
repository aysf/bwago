package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/aysf/bwago/internal/config"
	"github.com/aysf/bwago/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var templateDir = "./templates"

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData adds data for all template
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template renders template using html/template
func Template(rw http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template
	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		// create template cache
		tc, _ = CreateTemplateCache()
	}

	rt, ok := tc[tmpl]
	if !ok {
		log.Println("can't get template from cache")
		return errors.New("can't get tamplate from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := rt.Execute(buf, td)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(rw)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache creates template cache in a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	// tc is template cache
	tc := map[string]*template.Template{}

	// getting all page templates
	pages, err := filepath.Glob(templateDir + "/*.page.tmpl")
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
		layouts, err := filepath.Glob(templateDir + "/*.layout.tmpl")
		if err != nil {
			return tc, err
		}

		if len(layouts) > 0 {

			ts, err = ts.ParseGlob(templateDir + "/*.layout.tmpl")
			if err != nil {
				return tc, err
			}
		}

		tc[name] = ts
	}

	return tc, nil
}
