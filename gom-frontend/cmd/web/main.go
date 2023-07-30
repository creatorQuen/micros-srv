package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "tester-page.gohtml")
	})

	fmt.Println("Start gom front end service on port 5555")
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, templ string) {
	partials := []string{
		"./cmd/web/templates/base-layout.gohtml",
		"./cmd/web/templates/header-partial.gohtml",
		"./cmd/web/templates/footer-partial.gohtml",
	}

	var templates []string
	templates = append(templates, fmt.Sprintf("./cmd/web/templates/%s", templ))

	for _, part := range partials {
		templates = append(templates, part)
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
