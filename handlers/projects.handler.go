package handlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"nicholas-deary/config"
	"nicholas-deary/database"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type ProjectsData struct {
	Page        string
	Projects    []database.Project
	CurrentPage int
	TotalPages  int
}

type ProjectData struct {
	Page        string
	Project     *database.Project
	ArticleHTML template.HTML
}

func ProjectsHandler(w http.ResponseWriter, r *http.Request, t *template.Template, d *database.SQLiteDatabase) {
	values := r.URL.Query()
	page := 1
	search := values.Get("search")
	types := []string{}
	languages := []string{}
	technologies := []string{}

	if values.Has("page") {
		tmpPage, err := strconv.Atoi(values.Get("page"))
		if err == nil {
			page = tmpPage
		}
	}

	if values.Has("types") {
		types = strings.Split(values.Get("types"), ",")
	}
	if values.Has("languages") {
		languages = strings.Split(values.Get("languages"), ",")
	}
	if values.Has("technologies") {
		technologies = strings.Split(values.Get("technologies"), ",")
	}

	p, err := d.GetProjects(page, search, types, languages, technologies)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(p)

	err = t.ExecuteTemplate(w, "projects.gohtml", ProjectsData{
		Page: "/projects", 
		Projects: p, 
		CurrentPage: page, 
		TotalPages: int(math.Ceil(float64(d.GetProjectCount()) / 10)),
	})
	if err != nil {
		log.Print(err)
		return
	}
}

func ProjectHandler(w http.ResponseWriter, r *http.Request, c *config.Config, t *template.Template, d *database.SQLiteDatabase) {
	v := mux.Vars(r)
	name, err := url.QueryUnescape(v["name"])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(name) == 0 {
		log.Print("Empty name. Handlers set wrong.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p, err := d.GetProjectByName(name)
	if err != nil {
		log.Print(err)
		return
	}

	unsafe := blackfriday.Run(p.File)
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	html := policy.SanitizeBytes(unsafe)

	err = t.ExecuteTemplate(w, "project.gohtml", ProjectData{Page: "/projects", Project: p, ArticleHTML: template.HTML(html)})

	if err != nil {
		log.Print(err)
		return
	}
}

func PostProjectHandler(w http.ResponseWriter, r *http.Request, d *database.SQLiteDatabase) {
	var body struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Date        string `json:"date"`
		Type        string `json:"type"`
		Image       []byte `json:"image"`

		Languages    []string `json:"languages"`
		Technologies []string `json:"technologies"`

		File []byte `json:"file"`
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(b, &body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p := database.Project{
		ID:          body.ID,
		Name:        body.Name,
		Description: body.Description,
		Date:        body.Date,
		Type:        body.Type,

		Languages:    body.Languages,
		Technologies: body.Technologies,

		File: body.File,
	}

	project, err := d.CreateProject(p, body.Image)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	log.Print("New post submitted: " + project.Name)
	return
}

func getHTMLRenderer() {

}
