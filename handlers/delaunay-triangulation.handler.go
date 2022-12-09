package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type DelaunayTriangulationData struct {
	Page string
}

func DelaunayTriangulationHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	err := t.ExecuteTemplate(w, "delaunay-triangulation.gohtml", DelaunayTriangulationData{Page: "/projects"})

	if err != nil {
		log.Print(err)
		return
	}
}