package handlers

import (
	"net/http"
	"net/http/httptest"
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
	data               []postData
	expectedStatusCode int
}{
	{"home", "GET", "/", []postData{}, http.StatusOK},
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

		}

	}
}
