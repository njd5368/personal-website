package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"nicholas-deary/config"
	"nicholas-deary/database"
	"strconv"

	"github.com/gorilla/mux"
)

func GetImageHandler(w http.ResponseWriter, r *http.Request, d *database.SQLiteDatabase) {
	v := mux.Vars(r)
	idString, exists := v["id"]
	if !exists {
		log.Print("No image ID specified.")
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	image, err := d.GetImageByID(id)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(image)
	return
}

func PostImageHandler(w http.ResponseWriter, r *http.Request, d *database.SQLiteDatabase, c *config.Config) {
	i, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := d.CreateImage(i)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", c.Scheme + "://" + c.Host + "/image/" + strconv.FormatInt(id, 10))
	return
}