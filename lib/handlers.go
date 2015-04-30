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
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
	if status == http.StatusSeeOther {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HandleBase(context *AppContext, w http.ResponseWriter, r *http.Request) (retVal int, err error) {
	context.Log.Println("Inside the base handler.")
	switch r.Method {
	case "GET":
		code := r.URL.Query().Get("code")
		if instaError := r.URL.Query().Get("error"); len(code) != 0 {
			context.Log.Println("Send a post request to Instgaram now")
			return PerformPostReqeust(context, w, r, code)
		} else if instaError == "access_denied" && r.URL.Query().Get("error_reason") == "user_denied" {
			fmt.Fprintf(w, "Oops, something went wrong.")
			retVal = http.StatusInternalServerError
			err = errors.New("Seems you didn't allow that to happen")
		} else {
			t, err := template.ParseGlob("templates/*.html")
			if err != nil {
				context.Log.Println("Unable to parse the template. Error is: ", err)
				return http.StatusInternalServerError, err
			}
			t.Execute(w, nil)
			retVal = http.StatusFound
			err = nil
		}
	default:
		context.Log.Println("Seems the user didn't do a get request. Request type was: ", r.Method)
		retVal = http.StatusInternalServerError
		err = errors.New("Method not supported.")
	}
	return
}

func AuthroizeHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (revVal int, err error) {
	context.Log.Println("Inside the authorization handler.")
	return AuthenticateUser(context, w, r)
}
