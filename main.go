package main

import (
	"github.com/gotailwindcss/tailwind/twembed"
	"github.com/gotailwindcss/tailwind/twhandler"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var addr = ":3000"
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("site/templates/*"))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "site/images/favicon.ico")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)

	if err != nil {
		log.Println(err.Error())
	}
}

func main() {
	r := mux.NewRouter()
	r.Handle("/site/css/{css-file}", twhandler.New(http.Dir("./site/css"), "/site/css", twembed.New()))
	r.Handle("/site/images/{image}", http.StripPrefix("/site/images", http.FileServer(http.Dir("./site/images"))))
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.HandleFunc("/", indexHandler)

	log.Println("Starting server....")
	log.Fatal(http.ListenAndServe(addr, r))
}
