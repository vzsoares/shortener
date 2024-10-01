package main

import (
	"encoding/json"
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
	I18nStr string
	Stage   string
}

var ConstsMap = utils.ConstsMap{
	"API_BASE_URL":  fmt.Sprintf("%v/%v", os.Getenv("API_BASE_URL"), "/shortener/v1/public-api"),
	"SITE_BASE_URL": os.Getenv("FRONT_BASE_URL"),
}

var DefaultLang = "en"

func main() {
	stage := os.Getenv("STAGE")
	consts := utils.NewConsts(stage, ConstsMap)

	var i18n map[string]any
	fileBytes, _ := os.ReadFile("./i18n.json")
	err := json.Unmarshal(fileBytes, &i18n)

	funcT := func(s string) string {
		stps := strings.Split(s, ".")
		stps = append([]string{DefaultLang}, stps...)
		r := Walk(stps, i18n)
		return r
	}

	funcMap := map[string]interface{}{"T": funcT}
	funcMap = template.FuncMap(funcMap)

	data := &Data{
		Palette: Palette{
			Primary: "#532B88",
		},
		Consts: Consts{
			API_BASE_URL:  consts.GetConst("API_BASE_URL"),
			SITE_BASE_URL: consts.GetConst("SITE_BASE_URL"),
		},
		I18nStr: string(fileBytes),
		Stage:   stage,
	}

	templates := template.New("master").Funcs(funcMap)
	template.Must(templates.ParseGlob("src/**/*.go.html"))

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

		template.Must(tmpl.ParseFiles(v))
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

func Walk(p []string, m map[string]any) string {
	l := len(p)
	el := p[0]
	if l == 1 {
		return m[el].(string)
	}
	return Walk(p[1:], m[el].(map[string]any))
}
