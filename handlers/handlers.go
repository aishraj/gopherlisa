package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aishraj/gopherlisa/auth"
)

func HandleBase(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside the base handler.")
	BaseHandler(w, r)
}

func AuthroizeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside the authorization handler.")
	AuthHandler(w, r)
}

func BaseHandler(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	if instaError := r.URL.Query().Get("error"); len(code) != 0 {
		log.Println("Send a post request to Instgaram now")
		fmt.Fprintf(w, "Hi there, we trying to verfiy you now. Hang on tight.")
		auth.PerformPostReqeust(w, r, code)
	} else if instaError == "access_denied" && r.URL.Query().Get("error_reason") == "user_denied" {
		fmt.Fprintf(w, "Oops, something went wrong.")
	} else {
		t, err := template.ParseGlob("templates/*.html")
		if err != nil {
			log.Fatal("Unable to parse the template")
		}
		t.Execute(w, nil)
	}

}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	auth.AuthenticateUser(w, r)
}
