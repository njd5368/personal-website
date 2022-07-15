package handlers

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"nicholas-deary/config"
	"nicholas-deary/database"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

type PostsData struct {
	Page        	string
	Search 			string
	Categories 		[]string
	Languages		[]string
	AllLanguages 	[]string
	Technologies	[]string
	AllTechnologies []string
	Posts    	[]database.Post
	CurrentPage 	int
	TotalPages  	int
	LastPage		string
	NextPage		string
}

type PostData struct {
	Page        string
	Post     *database.Post
	ArticleHTML template.HTML
}

func ProjectsHandler(w http.ResponseWriter, r *http.Request, t *template.Template, d *database.SQLiteDatabase) {
	ProjectsAndBlogHandler(w, r, t, d, database.Project)
}

func BlogHandler(w http.ResponseWriter, r *http.Request, t *template.Template, d *database.SQLiteDatabase) {
	ProjectsAndBlogHandler(w, r, t, d, database.Blog)
}

func ProjectsAndBlogHandler(w http.ResponseWriter, r *http.Request, t *template.Template, d *database.SQLiteDatabase, urlRoute database.Type) {
	values := r.URL.Query()
	page := 1
	search := values.Get("search")
	categories := []string{}
	languages := []string{}
	technologies := []string{}
	first := ""
	last := ""

	if values.Has("types") {
		categories = values["types"]
	}
	if values.Has("languages") {
		languages = values["languages"]
	}
	if values.Has("technologies") {
		technologies = values["technologies"]
	}

	allLanguages, err := d.GetAllLanguages()
	if err != nil {
		allLanguages = []string{}
	}

	allTechnologies, err := d.GetAllTechnologies()
	if err != nil {
		allTechnologies = []string{}
	}

	var totalPages int
	if urlRoute == database.Project {
		totalPages = int(math.Ceil(float64(d.GetProjectPostCount(search, categories, languages, technologies)) / 10))
	} else {
		totalPages = int(math.Ceil(float64(d.GetBlogPostCount(search, categories, languages, technologies)) / 10))
	}
	

	if values.Has("page") {
		tmpPage, err := strconv.Atoi(values.Get("page"))
		if err == nil {
			if tmpPage < 1 {
				page = 1
			} else if tmpPage > totalPages {
				page = totalPages
			} else {
				page = tmpPage
			}
		}
	}

	if values.Has("f") {
		first = values.Get("f")
	}
	if values.Has("l") {
		last = values.Get("l")
	}
	
	var p []database.Post
	var navpage string
	if urlRoute == database.Project {
		p, err = d.GetProjectPosts(page, search, categories, languages, technologies, first, last)
		navpage = "projects"

	} else {
		p, err = d.GetBlogPosts(page, search, categories, languages, technologies, first, last)
		navpage = "blog"
	}
	
	if err != nil {
		log.Print(err)
		return
	}

	var lastPage string
	if page != 1 {
		values.Del("l")
		values.Set("f", p[0].Date + " " + strconv.FormatInt(p[0].ID, 10))
		values.Set("page", strconv.Itoa(page - 1))
		lastPage = "/" + navpage + "?" + values.Encode()
	} else {
		lastPage = "#"
	}
	
	var nextPage string
	if page < totalPages {
		values.Del("f")
		values.Set("l", p[len(p) - 1].Date + " " + strconv.FormatInt(p[len(p) - 1].ID, 10))
		values.Set("page", strconv.Itoa(page + 1))
		nextPage = "/" + navpage + "?" + values.Encode()
	} else {
		nextPage = "#"
	}

	err = t.ExecuteTemplate(w, navpage + ".gohtml", PostsData{
		Page: 				"/" + navpage, 
		Posts: 				p,
		Search: 			search,
		Categories: 		categories,
		Languages: 			languages,
		AllLanguages: 		allLanguages,
		Technologies: 		technologies,
		AllTechnologies:	allTechnologies,
		CurrentPage: 		page, 
		TotalPages: 		totalPages,
		LastPage: 			lastPage,
		NextPage: 			nextPage,
	})
	if err != nil {
		log.Print(err)
		return
	}
}

func GetProjectPostHandler(w http.ResponseWriter, r *http.Request, c *config.Config, t *template.Template, d *database.SQLiteDatabase) {
	GetPostHandler(w, r, c, t, d, database.Project)
}

func GetBlogPostHandler(w http.ResponseWriter, r *http.Request, c *config.Config, t *template.Template, d *database.SQLiteDatabase) {
	GetPostHandler(w, r, c, t, d, database.Blog)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request, c *config.Config, t *template.Template, d *database.SQLiteDatabase, urlRoute database.Type) {
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

	var p *database.Post
	var navpage string
	if urlRoute == database.Project {
		p, err = d.GetProjectPostByName(name)
		navpage = "/projects"
	} else {
		p, err = d.GetBlogPostByName(name)
		navpage = "/blog"
	}
	if err != nil {
		log.Print(err)
		return
	}

	unsafe := blackfriday.Run(p.File)
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	html := policy.SanitizeBytes(unsafe)

	err = t.ExecuteTemplate(w, "post.gohtml", PostData{Page: navpage, Post: p, ArticleHTML: template.HTML(html)})

	if err != nil {
		log.Print(err)
		return
	}
}
