package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"nicholas-deary/handlers"
)

var addr = ":3000"
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("site/templates/util/*"))
	tpl = template.Must(tpl.ParseGlob("site/templates/pages/*"))
	tpl = template.Must(tpl.ParseGlob("site/templates/error-pages/*"))
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

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.NotFoundHandler(w, r, tpl)
	})

	log.Println("Starting server....")
	log.Fatal(http.ListenAndServe(addr, r))
}
