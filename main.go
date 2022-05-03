package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"nicholas-deary/config"
	"nicholas-deary/database"
	"nicholas-deary/handlers"
	"strconv"

	"github.com/gorilla/mux"
)

type Config struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Port   int    `json:"port"`

	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`

	Users []struct{
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"users"`
}

var t *template.Template

func init() {
	t = template.Must(template.ParseGlob("site/templates/util/*"))
	t = template.Must(t.ParseGlob("site/templates/pages/*"))
	t = template.Must(t.ParseGlob("site/templates/error-pages/*"))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "site/images/favicon.ico")
}

func main() {
	c, err := config.ReadConfigFile()
	if err != nil {
		log.Panic(err)
	}

	d := database.NewSQLiteDatabase(&sql.DB{}) //TODO: add db

	r := mux.NewRouter()
	a := r.NewRoute().Subrouter()

	r.Handle("/site/css/{css-file}", http.StripPrefix("/site/css", http.FileServer(http.Dir("./site/css"))))
	r.Handle("/site/images/{image}", http.StripPrefix("/site/images", http.FileServer(http.Dir("./site/images"))))
	r.HandleFunc("/favicon.ico", faviconHandler)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r, t)
	})
	r.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		handlers.AboutHandler(w, r, t)
	})
	r.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProjectsHandler(w, r, t, d)
	}).Methods("GET")
	r.HandleFunc("projects/{name}", func(w http.ResponseWriter, r *http.Request) {

	})
	r.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		handlers.BlogHandler(w, r, t)
	})
	r.HandleFunc("blog/{name}", func(w http.ResponseWriter, r *http.Request) {

	})

	r.HandleFunc("/image/{id}", func(w http.ResponseWriter, r *http.Request) {

	})

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.NotFoundHandler(w, r, t)
	})

	a.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostProjectHandler(w, r, d)
	}).Methods("POST")

	a.Use()

	if len(c.CertFile) == 0 || len(c.KeyFile) == 0 {
		log.Println("Starting insecure server " + c.Host + " on port " + strconv.Itoa(c.Port) + "....")
		log.Fatal(http.ListenAndServe(c.Host+":"+strconv.Itoa(c.Port), r))
	} else {
		log.Println("Starting secure server " + c.Host + " on port " + strconv.Itoa(c.Port) + "....")
		log.Fatal(http.ListenAndServeTLS(c.Host+":"+strconv.Itoa(c.Port), c.CertFile, c.KeyFile, r))	
	}
}
