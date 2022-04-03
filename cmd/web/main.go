package main

import (
	"log"
	"net/http"

	"github.com/aysf/bwago/pkg/config"
	"github.com/aysf/bwago/pkg/render"
)

var portNumber = ":8080"

func main() {

	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Panic("error loading template cache ", err)
	}

	app.TemplateCache = tc
	app.UseCache = true

	render.NewTemplates(&app)

	routes(&app)

	// http.HandleFunc("/", handlers.Home)
	// http.HandleFunc("/about", handlers.About)

	http.Handle("/src/", http.StripPrefix("/src", http.FileServer(http.Dir("./static"))))

	log.Println("application running on port 8080")
	// http.ListenAndServe(":8080", nil)

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	server.ListenAndServe()
	log.Fatal(err)
}
