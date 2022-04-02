package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// RenderTemplate renders template using html/template
func RenderTemplate(rw http.ResponseWriter, tmpl string) error {

	tc, _ := CreateTemplateCache()

	rt, ok := tc[tmpl]
	if !ok {
		log.Panic("template does not exist")
	}

	buf := new(bytes.Buffer)
	rt.Execute(buf, nil)
	_, err := buf.WriteTo(rw)
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
