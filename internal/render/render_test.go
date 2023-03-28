package render

import (
	"net/http"
	"testing"

	"github.com/aysf/bwago/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	res := AddDefaultData(&td, r)
	if res.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {

	templateDir = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, _ := getSession()

	ww := new(myWriter)

	err = Template(ww, r, "about.page.gohtml", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(ww, r, "non-exist.page.gohtml", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template does not exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()

	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {

	templateDir = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error("error creating template cache")
	}
}
