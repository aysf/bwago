package main

import (
	"net/http"
	"os"
	"testing"
)

// TestMain run before test run
func TestMain(m *testing.M) {

	os.Exit(m.Run())

}

type MyHandler struct{}

func (mh *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
