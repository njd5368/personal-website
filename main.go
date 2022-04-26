package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"nicholas-deary/handlers"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Config struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Port   int    `json:"port"`

	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("site/templates/util/*"))
	tpl = template.Must(tpl.ParseGlob("site/templates/pages/*"))
	tpl = template.Must(tpl.ParseGlob("site/templates/error-pages/*"))
}

func readConfigFile() Config {
	h, err := os.UserHomeDir()

	if err != nil {
		log.Panic("Problem getting home directory: " + err.Error())
	}

	b, err := ioutil.ReadFile(h + "/.config/personal-website/config.json")

	if errors.Is(err, os.ErrNotExist) {
		log.Panic("Config file does not exist at ~/.config/personal-website/config.json\n" + err.Error())
	}
	if err != nil {
		log.Panic("Problem opening ~/.config/personal-website/config.json\n" + err.Error())
	}

	if err != nil {
		log.Panic("Problem reading config file\n" + err.Error())
	}

	var c Config
	err = json.Unmarshal(b, &c)

	if err != nil {
		log.Panic("Problem unmarshaling data from config\n" + err.Error())
	}

	return c
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "site/images/favicon.ico")
}

func main() {
	r := mux.NewRouter()
	r.Handle("/site/css/{css-file}", http.StripPrefix("/site/css", http.FileServer(http.Dir("./site/css"))))
	r.Handle("/site/images/{image}", http.StripPrefix("/site/images", http.FileServer(http.Dir("./site/images"))))
	r.HandleFunc("/favicon.ico", faviconHandler)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r, tpl)
	})
	r.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		handlers.AboutHandler(w, r, tpl)
	})
	r.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProjectsHandler(w, r, tpl)
	})
	r.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		handlers.BlogHandler(w, r, tpl)
	})

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.NotFoundHandler(w, r, tpl)
	})

	c := readConfigFile()

	if len(c.CertFile) == 0 || len(c.KeyFile) == 0 {
		log.Println("Starting insecure server " + c.Host + " on port " + strconv.Itoa(c.Port) + "....")
		log.Fatal(http.ListenAndServe(c.Host+":"+strconv.Itoa(c.Port), r))
	} else {
		log.Println("Starting secure server " + c.Host + " on port " + strconv.Itoa(c.Port) + "....")
		log.Fatal(http.ListenAndServeTLS(c.Host+":"+strconv.Itoa(c.Port), c.CertFile, c.KeyFile, r))	
	}
}
