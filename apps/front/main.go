package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	templates, err := template.ParseGlob("src/**/*.go.html")

	if err != nil {
		log.Fatal(err)
	}
	w, err := os.Create("dist/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = templates.ExecuteTemplate(w, "index.go.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}
