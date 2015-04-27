package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aishraj/gopherlisa/auth"
)

func BaseHandler(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	if instaError := r.URL.Query().Get("error"); len(code) != 0 {
		//now we have the code from instagram
		fmt.Fprintf(w, "Hi there, we trying to verfiy you now. Hang on tight.")
		performPostReqeust(code)

	} else if instaError == "access_denied" && r.URL.Query().Get("error_reason") == "user_denied" {
		fmt.Fprintf(w, "Oops, something went wrong.")
	} else {
		t, err := template.ParseFiles("templates/login.html")
		if err != nil {
			log.Fatal("Unable to parse the template")
		}
		t.Execute(w, nil)
	}

	//TODO: add more code to handle the callback

}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	//This one should now direct to instagram.
	auth.AuthenticateUser(w, r)
}
