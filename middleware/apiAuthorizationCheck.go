package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"
	"nicholas-deary/config"
	"strings"

	"golang.org/x/crypto/sha3"
    "golang.org/x/exp/slices"
)

type APIAuthorization struct {
	Config *config.Config
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

        hash := sha3.New256()
		hash.Write(authorizationLoginSplit[1])
		password := hex.EncodeToString(hash.Sum(nil))

        i := slices.IndexFunc(a.Config.Users, func(u config.User) bool {
            return u.Username == username
        })
        if i == -1 || password != a.Config.Users[i].Password {
            log.Printf("Authorization attempt was made by %s.", username)
            w.WriteHeader(300)
            w.Write([]byte("Unauthorized."))
            return
        }
        log.Printf("%s is making an API call.", username)

        next.ServeHTTP(w, r)
    })
}