package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	tmpl, err := template.ParseFiles("src/pages/index.go.html")
	if err != nil {
		log.Fatal(err)
	}
	w, err := os.Create("dist/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}
