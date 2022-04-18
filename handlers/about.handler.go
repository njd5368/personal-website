package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type AboutData struct {
	Page string
}

func AboutHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	err := t.ExecuteTemplate(w, "about.gohtml", AboutData{Page: "/about"})

	if err != nil {
		log.Print(err)
		return
	}
}