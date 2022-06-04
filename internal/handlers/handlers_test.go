package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var testCase = []struct {
	name               string
	method             string
	path               string
	params             []postData
	expectedStatusCode int
}{
	{"home", "GET", "/", []postData{}, http.StatusOK},
	{"about", "GET", "/about", []postData{}, http.StatusOK},
	{"contact", "GET", "/contact", []postData{}, http.StatusOK},
	{"generals-quarters", "GET", "/generals-quarters", []postData{}, http.StatusOK},
	{"majors-suite", "GET", "/majors-suite", []postData{}, http.StatusOK},
	{"search-availability", "GET", "/search-availability", []postData{}, http.StatusOK},
	{"make-reservation", "GET", "/make-reservation", []postData{}, http.StatusOK},
	{"reservation-summary", "GET", "/reservation-summary", []postData{}, http.StatusOK},
	{"search-availability", "POST", "/search-availability", []postData{
		{key: "start", value: "01-02-2021"},
		{key: "end", value: "01-03-2021"},
	}, http.StatusOK},
	{"search-availability-json", "POST", "/search-availability-json", []postData{
		{key: "start", value: "01-02-2021"},
		{key: "end", value: "01-03-2021"},
	}, http.StatusOK},
	{"make-reservation", "POST", "/make-reservation", []postData{
		{key: "first_name", value: "Tra"},
		{key: "last_name", value: "Guy"},
		{key: "phone", value: "555-555"},
		{key: "email", value: "ansuf@gmail.com"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoute()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range testCase {

		if e.method == "GET" {

			resp, err := ts.Client().Get(ts.URL + e.path)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected status code %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		} else {
			values := url.Values{}

			for _, x := range e.params {
				values.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.path, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}

	}
}
