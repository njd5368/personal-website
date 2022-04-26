package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type BlogData struct {
	Page string
}

func BlogHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	err := t.ExecuteTemplate(w, "blog.gohtml", ProjectsData{Page: "/blog"})

	if err != nil {
		log.Print(err)
		return
	}
}