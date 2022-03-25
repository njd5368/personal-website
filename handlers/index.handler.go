package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	err := t.ExecuteTemplate(w, "index.gohtml", nil)

	if err != nil {
		log.Println(err.Error())
		return
	}
}
