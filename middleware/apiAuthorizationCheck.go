package middleware

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"nicholas-deary/config"
	"nicholas-deary/database"
	"strings"
)

type APIAuthorization struct {
	Config   *config.Config
	Database *database.SQLiteDatabase
}

func (a APIAuthorization) CheckUserAuthorziation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("API call made from " + r.RemoteAddr)

		authorization := r.Header.Get("Authorization")

		authorizationList := strings.Split(authorization, " ")

		if authorizationList[0] != "Basic" {
			log.Print("Authorization attempt with bad authorization method.")
			w.WriteHeader(300)
			w.Write([]byte("Authorization method not recognized."))
			return
		}

		authorizationLogin, err := base64.StdEncoding.DecodeString(authorizationList[1])
		if err != nil {
			log.Print("Authorization attempt without base64 encoding.")
			w.WriteHeader(300)
			w.Write([]byte("Base64 encoding incorrect."))
			return
		}

		authorizationLoginSplit := bytes.Split(authorizationLogin, []byte(":"))

		username := string(authorizationLoginSplit[0])

		password := string(authorizationLoginSplit[1])

		if !a.Database.IsUser(username, password) {
			log.Printf("Authorization attempt was made by %s.", username)
			w.WriteHeader(300)
			w.Write([]byte("Unauthorized."))
			return
		}
		log.Printf("%s is making an API call.", username)

		next.ServeHTTP(w, r)
	})
}

