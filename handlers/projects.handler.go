package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type ProjectsData struct {
	Page string
}

func ProjectsHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	err := t.ExecuteTemplate(w, "projects.gohtml", ProjectsData{Page: "/projects"})

	if err != nil {
		log.Print(err)
		return
	}
}