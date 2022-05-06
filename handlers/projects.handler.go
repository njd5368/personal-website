package handlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"nicholas-deary/database"
)

type ProjectsData struct {
	Page string
	Projects []database.Project
}

func ProjectsHandler(w http.ResponseWriter, r *http.Request, t *template.Template, d *database.SQLiteDatabase) {
	p, err := d.GetProjects()
	if err != nil {
		log.Print(err)
		return
	}

	err = t.ExecuteTemplate(w, "projects.gohtml", ProjectsData{Page: "/projects", Projects: p})

	if err != nil {
		log.Print(err)
		return
	}
}

func PostProjectHandler(w http.ResponseWriter, r *http.Request, d *database.SQLiteDatabase) {
	var body struct{
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Date        string `json:"date"`
		Type        string `json:"type"`
		Image       []byte  `json:"image"`

		Languages    []string `json:"languages"`
		Technologies []string `json:"technologies"`

		File []byte `json:"file"`
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}

	p := database.Project{
		ID: body.ID,
		Name: body.Name,
		Description: body.Description,
		Date: body.Date,
		Type: body.Type,

		Languages: body.Languages,
		Technologies: body.Technologies,

		File: body.File,
	}

	_, err = d.CreateProject(p, body.Image)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
}
