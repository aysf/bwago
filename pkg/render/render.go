package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/aysf/bwago/pkg/config"
	"github.com/aysf/bwago/pkg/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders template using html/template
func RenderTemplate(rw http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

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
		log.Panic("template does not exist")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	err := rt.Execute(buf, td)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(rw)
	if err != nil {
		return err
	}

	return nil
}

// CreateTemplateCache creates template cache in a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	// tc is template cache
	tc := map[string]*template.Template{}

	// getting all page templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
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
		layouts, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return tc, err
		}

		if len(layouts) > 0 {

			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return tc, err
			}
		}

		tc[name] = ts
	}

	return tc, nil
}
