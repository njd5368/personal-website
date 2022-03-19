package main

import (
	"html/template"
	"io/ioutil"
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

func imageHandler(w http.ResponseWriter, r *http.Request) {

	fileBytes, err := ioutil.ReadFile("site/images/signal.jpeg")
	if err != nil {
		panic(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(fileBytes); err != nil {
		panic(err.Error())
	}
	return
}

func main() {
	r := mux.NewRouter()

	r.Handle("/site/images/", http.StripPrefix("/site/images/", http.FileServer(http.Dir("/site/images"))))
	r.HandleFunc("/site/images/", imageHandler)
	r.HandleFunc("/", indexHandler)

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(addr, r))
}
