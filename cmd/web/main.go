package main

import (
	"log"
	"net/http"

	"github.com/aysf/bwago/pkg/handlers"
)

func main() {

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	http.Handle("/src/", http.StripPrefix("/src", http.FileServer(http.Dir("./static"))))

	log.Println("application running on port 8080")
	http.ListenAndServe(":8080", nil)
}
