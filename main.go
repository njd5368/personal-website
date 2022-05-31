package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"nicholas-deary/config"
	"nicholas-deary/database"
	"nicholas-deary/handlers"
	"nicholas-deary/middleware"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const databaseFile = "blog.db"
var t *template.Template
var c *config.Config

func init() {
	var err error
	c, err = config.ReadConfigFile()
	if err != nil {
		log.Panic(err)
	}

	t = template.New("").Funcs(template.FuncMap{
        "imageFromID": func(i int64) string {
			return c.Scheme + "://" + c.Host + ":" + strconv.Itoa(c.Port) + "/image/" + strconv.FormatInt(i, 10)
		},
		"fmtDate": func(d string) string {
			t, err := time.Parse("2006-01-02", d)
			if err != nil {
				return d
			}
			return t.Format("Jan 02, 2006")
		},
		"projectColor": func(p string) string {
			colors := map[string]string{
				"personal": "#F28FAD",
				"hackathon": "#F8BD96",
				"professional": "#ABE9B3",
				"academic": "#96CDFB",
			}
			return colors[strings.ToLower(p)]
		},
		"getPageNumbers": func(c int, t int) []string {
			if c < 1 || t < 1 {
				return []string{}
			}
			result := []string{}
			if t < 12 {
				for i := 1; i <= t; i++ {
					result = append(result, strconv.Itoa(i))
				}
			} else if c < 7 {
				for i := 1; i <= 9; i ++ {
					result = append(result, strconv.Itoa(i))
				}
				result = append(result, "...", strconv.Itoa(t))
			} else if t - c < 6 {
				result = append(result, "1", "...")
				for i := t - 8; i <= t; i++ {
					result = append(result, strconv.Itoa(i))
				}
			} else {
				result = append(result, "1", "...")
				for i := c - 3; i <= c + 3; i++ {
					result = append(result, strconv.Itoa(i))
				}
				result = append(result, "...", strconv.Itoa(t))
			}
			return result
		},
		"intToString": func(i int) string {
			return strconv.Itoa(i)
		},
    })
	t = template.Must(t.ParseGlob("site/templates/util/*"))
	t = template.Must(t.ParseGlob("site/templates/pages/*"))
	t = template.Must(t.ParseGlob("site/templates/error-pages/*"))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "site/images/favicon.ico")
}

func main() {
	log.Print("Configuring relational database...")
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Panic(err)
	}

	d := database.NewSQLiteDatabase(db)
	err = d.StartDatabase()
	if err != nil {
		log.Panic(err)
	}

	log.Print("Setting up routes...")
	r := mux.NewRouter()
	a := r.NewRoute().Subrouter()

	r.Handle("/site/css/{css-file}", http.StripPrefix("/site/css", http.FileServer(http.Dir("./site/css"))))
	r.Handle("/site/js/{js-file}", http.StripPrefix("/site/js", http.FileServer(http.Dir("./site/js"))))
	r.Handle("/site/images/{image}", http.StripPrefix("/site/images", http.FileServer(http.Dir("./site/images"))))
	r.HandleFunc("/favicon.ico", faviconHandler)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r, t)
	}).Methods("GET")
	r.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		handlers.AboutHandler(w, r, t)
	}).Methods("GET")
	r.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProjectsHandler(w, r, t, d)
	}).Methods("GET")
	r.HandleFunc("/projects/{name}", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProjectHandler(w, r, c, t, d)
	}).Methods("GET")
	r.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		handlers.BlogHandler(w, r, t)
	}).Methods("GET")
	r.HandleFunc("/blog/{name}", func(w http.ResponseWriter, r *http.Request) {

	}).Methods("GET")

	r.HandleFunc("/image/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetImageHandler(w, r, d)
	}).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.NotFoundHandler(w, r, t)
	})

	a.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostProjectHandler(w, r, d)
	}).Methods("POST")
	a.HandleFunc("/api/image", func (w http.ResponseWriter, r *http.Request)  {
		handlers.PostImageHandler(w, r, d, c)
	}).Methods("POST")

	a.Use(middleware.APIAuthorization{Config: c}.CheckUserAuthorziation)

	if len(c.CertFile) == 0 || len(c.KeyFile) == 0 {
		log.Println("Starting insecure server " + c.Host + " on port " + strconv.Itoa(c.Port) + "....")
		log.Fatal(http.ListenAndServe(c.Host+":"+strconv.Itoa(c.Port), r))
	} else {
		log.Println("Starting secure server " + c.Host + " on port " + strconv.Itoa(c.Port) + "....")
		log.Fatal(http.ListenAndServeTLS(c.Host+":"+strconv.Itoa(c.Port), c.CertFile, c.KeyFile, r))	
	}
}
