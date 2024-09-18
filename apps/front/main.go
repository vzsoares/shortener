package main

//go:generate make css-build

import (
	"fmt"
	"html/template"
	"io"
	"libs/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Palette struct {
	Primary string
}
type Consts struct {
	API_BASE_URL  string
	SITE_BASE_URL string
}

type Data struct {
	Palette Palette
	Consts  Consts
}

var DevConsts = utils.ConstsMap{
	"API_BASE_URL":  "https://api-dev.zenhalab.com/shortener/v1/public-api",
	"SITE_BASE_URL": "https://s-dev.zenhalab.com",
}
var ProdConsts = utils.ConstsMap{
	"API_BASE_URL":  "https://api.zenhalab.com/shortener/v1/public-api",
	"SITE_BASE_URL": "https://s.zenhalab.com",
}

func main() {
	consts := utils.NewConsts(os.Getenv("STAGE"), ProdConsts, DevConsts)

	data := &Data{
		Palette: Palette{
			Primary: "#532B88",
		},
		Consts: Consts{
			API_BASE_URL:  consts.GetConst("API_BASE_URL"),
			SITE_BASE_URL: consts.GetConst("SITE_BASE_URL"),
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

	pages, err := filepath.Glob("src/pages/*")
	for _, v := range pages {
		fname := filepath.Base(v)
		bname := strings.Split(fname, ".")[0]

		w, err := os.Create(fmt.Sprintf("dist/%v.html", bname))
		if err != nil {
			log.Fatal(err)
		}
		// clone and reparse file fixes redefined template names
		// https://stackoverflow.com/a/69244593/16923160
		tmpl, err := templates.Clone()
		if err != nil {
			log.Fatal(err)
		}

		tmpl, err = tmpl.ParseFiles(v)
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.ExecuteTemplate(w, fname, data)
		if err != nil {
			log.Fatal(err)
		}
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
