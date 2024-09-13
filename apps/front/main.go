package main

import (
	"html/template"
	"io"
	"log"
	"os"
)

type Palette struct {
	Primary string
}

type Data struct {
	Palette Palette
}

func main() {
	data := &Data{
		Palette: Palette{
			Primary: "#532B88",
		},
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
	err = Copy("./src/assets/favicon.ico", "./dist/assets/favicon.ico")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Copy(srcpath, dstpath string) (err error) {
	r, err := os.Open(srcpath)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dstpath)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	return err
}
