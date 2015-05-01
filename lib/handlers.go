package lib

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type AppContext struct {
	Log          *log.Logger
	SessionStore *SessionManager
}

type tinyUser struct {
	DisplayName string
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
	session := context.SessionStore.SessionStart(w, r)
	context.Log.Println("Got session object.")
	createtime := session.Get("createtime")
	context.Log.Println("The create time of the session is: ", createtime)
	if createtime == nil {
		session.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		context.SessionStore.SessionDestroy(session.SessionID())
		session = context.SessionStore.SessionStart(w, r)
	}
	context.Log.Println("Now checking the method.")
	switch r.Method {
	case "GET":
		context.Log.Println("It is a GET method.")
		code := r.URL.Query().Get("code")
		if instaError := r.URL.Query().Get("error"); len(code) != 0 {
			context.Log.Println("Send a post request to Instgaram now")
			return PerformPostReqeust(context, w, r, code)
		} else if instaError == "access_denied" && r.URL.Query().Get("error_reason") == "user_denied" {
			fmt.Fprintf(w, "Oops, something went wrong.")
			retVal = http.StatusInternalServerError
			err = errors.New("Seems you didn't allow that to happen")
		} else {
			context.Log.Println("Now trying to get the value of a field from the session")
			displayUser := session.Get("user")
			if displayUser == nil {
				context.Log.Println("The value was not set. Setting it to guest now")
				session.Set("user", "Guest")
			} else {
				context.Log.Println("The value was set, displaying it now.")
				session.Set("user", displayUser)
			}
			context.Log.Println("Now trying to render the template.")
			t, err := template.ParseGlob("templates/*.html")
			if err != nil {
				context.Log.Println("Unable to parse the template. Error is: ", err)
				return http.StatusInternalServerError, err
			}
			w.Header().Set("Content-Type", "text/html")

			context.Log.Println("Display user is: ", session.Get("user"))
			if str, ok := session.Get("user").(string); ok {
				userInfo := tinyUser{str}
				t.Execute(w, userInfo)
				retVal = http.StatusFound
				err = nil
			} else {
				context.Log.Printf("Could not cast the user display name to string")
				retVal = http.StatusInternalServerError
				err = errors.New("Unable to set username in template")
			}
		}
	default:
		context.Log.Println("Seems the user didn't do a get request. Request type was: ", r.Method)
		retVal = http.StatusInternalServerError
		err = errors.New("Method not supported.")
	}
	context.Log.Printf("Returning with values Status  %v and error %v ", retVal, err)
	return
}

func AuthroizeHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (revVal int, err error) {
	session := context.SessionStore.SessionStart(w, r)
	authToken := session.Get("access_token")
	if authToken == nil {
		context.Log.Println("Inside the authorization handler.")
		return AuthenticateUser(context, w, r)
	} else {
		context.Log.Println("Already authenticated, and therefore returning from here.")
		return http.StatusSeeOther, nil
	}
}
