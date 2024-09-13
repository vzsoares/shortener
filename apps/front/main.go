package main

//go:generate make css-build

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	os.Mkdir("./dist", os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}

	w, err := os.Create("dist/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = templates.ExecuteTemplate(w, "index.go.html", data)
	if err != nil {
		log.Fatal(err)
	}

	// copy assets
	assets, err := filepath.Glob("src/assets/*")
	if len(assets) < 1 {
		return
	}

	os.Mkdir("./dist/assets", os.ModePerm)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, v := range assets {
		fname := filepath.Base(v)
		err = Copy(
			fmt.Sprintf("./src/assets/%v", fname),
			fmt.Sprintf("./dist/assets/%v", fname),
		)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// make css
	cssBuildCmd := exec.Command("make", "css-build")
	err = cssBuildCmd.Run()
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
