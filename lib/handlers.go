package lib

import (
	"errors"
	"fmt"
	"github.com/aishraj/gopherlisa/lib/sessions"
	"html/template"
	"log"
	"net/http"
)

type AppContext struct {
	Log          *log.Logger
	SessionStore *sessions.Manager
}

type Handler struct {
	*AppContext
	HandlerFunc func(*AppContext, http.ResponseWriter, *http.Request) (int, error)
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Updated to pass ah.appContext as a parameter to our handler type.
	status, err := handler.HandlerFunc(handler.AppContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			// And if we wanted a friendlier error page, we can
			// now leverage our context instance - e.g.
			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		case http.StatusSeeOther:
			http.Redirect(w, r, "/", http.StatusSeeOther)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

func HandleBase(context *AppContext, w http.ResponseWriter, r *http.Request) (retVal int, err error) {
	log.Println("Inside the base handler.")
	switch r.Method {
	case "GET":
		code := r.URL.Query().Get("code")
		if instaError := r.URL.Query().Get("error"); len(code) != 0 {
			log.Println("Send a post request to Instgaram now")
			return PerformPostReqeust(context, w, r, code)
		} else if instaError == "access_denied" && r.URL.Query().Get("error_reason") == "user_denied" {
			fmt.Fprintf(w, "Oops, something went wrong.")
			retVal = 422
			err = errors.New("Seems you didn't allow that to happen")
		} else {
			loadAndParseTemplate(w, r)
			retVal = 200
			err = nil
		}
	default:
		context.Log.Println("Seems the user didn't do a get request. Request type was: ", r.Method)
		retVal = 500
		err = errors.New("Method not supported.")
	}
	return
}

func AuthroizeHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (revVal int, err error) {
	log.Println("Inside the authorization handler.")
	return AuthenticateUser(context, w, r)
}

func loadAndParseTemplate(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseGlob("../templates/*.html")
	if err != nil {
		log.Fatal("Unable to parse the template")
	}
	t.Execute(w, nil)
}
