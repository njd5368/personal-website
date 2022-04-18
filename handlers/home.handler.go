package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type HomeData struct {
	Page string
}

func HomeHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	err := t.ExecuteTemplate(w, "home.gohtml", HomeData{Page: "/"})

	if err != nil {
		log.Print(err)
		return
	}
}
