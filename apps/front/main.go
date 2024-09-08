package main

import (
	"html/template"
	"log"
	"os"
)

type Palette struct {
	Primary string
}

type Data struct {
	Palette Palette
}

var Coolors = &Palette{
	Primary: "#532B88",
}

func main() {
	data := &Data{
		Palette: *Coolors,
	}

	templates, err := template.ParseGlob("src/**/*.go.html")
	if err != nil {
		log.Fatal(err)
	}
	w, err := os.Create("dist/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = templates.ExecuteTemplate(w, "index.go.html", data)
	if err != nil {
		log.Fatal(err)
	}
	// TODO mv assets folder to dist
}
