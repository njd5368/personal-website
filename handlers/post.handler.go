package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"nicholas-deary/config"
	"nicholas-deary/database"
	"strconv"

	"golang.org/x/exp/slices"
)

func PostPostHandler(w http.ResponseWriter, r *http.Request, d *database.SQLiteDatabase) {
	var body struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Date        string `json:"date"`
		Category    string `json:"category"`
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

	p := database.Post{
		ID:          body.ID,
		Name:        body.Name,
		Type: 		   database.StringToType(body.Type),
		Description: body.Description,
		Date:        body.Date,
		Category:    body.Category,

		Languages:    body.Languages,
		Technologies: body.Technologies,

		File: body.File,
	}
	
	var categories []string
	if p.Type == database.Project {
		categories = []string{"Personal", "Hackathon", "Professional", "Academic"}
	} else {
		categories = []string{"Development", "Linux", "Update", "Other"}
	}
	if !slices.Contains(categories, p.Category) {
		log.Print("Post type incorrect.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	project, err := d.CreatePost(p, body.Image)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	log.Print("New post submitted: " + project.Name)
}

func LatestPostHandler(w http.ResponseWriter, r *http.Request, c *config.Config, d *database.SQLiteDatabase) {
	name, err := d.GetLatestBlogPostName()
	if err != nil {
		log.Print(err)
		name = ""
	}
	
	link := c.Scheme + "://" + c.Host + ":" + strconv.Itoa(c.Port) + "/blog"
	if len(name) != 0 {
		link += "/" + url.QueryEscape(name)
	}
	
	http.Redirect(w, r, link, http.StatusSeeOther)
}
