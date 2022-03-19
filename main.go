package main

import (
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)

	if err != nil {
		log.Println(err.Error())
	}
}

func main() {
	r := mux.NewRouter()
	r.Handle("/site/images/{image}", http.StripPrefix("/site/images", http.FileServer(http.Dir("./site/images"))))
	r.HandleFunc("/", indexHandler)

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(addr, r))
}
