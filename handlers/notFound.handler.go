package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	err := t.ExecuteTemplate(w, "notFound.gohtml", nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
