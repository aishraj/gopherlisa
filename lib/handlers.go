package lib

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
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

func validateAndStartSession(context *AppContext, w http.ResponseWriter, r *http.Request) Session {
	session := context.SessionStore.SessionStart(w, r)
	createtime := session.Get("createtime")
	if createtime == nil {
		session.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		context.Log.Println("Session has expired, starting new session.")
		context.SessionStore.SessionDestroy(session.SessionID())
		session = context.SessionStore.SessionStart(w, r)
	}
	return session
}

func BaseHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (retVal int, err error) {
	session := validateAndStartSession(context, w, r)
	switch r.Method {
	case "GET":
		code := r.URL.Query().Get("code")
		instaError := r.URL.Query().Get("error")
		if instaError == "access_denied" && r.URL.Query().Get("error_reason") == "user_denied" {
			context.Log.Println("User denied permission for access.")
			retVal = http.StatusUnauthorized
			err = errors.New("Seems you didn't allow that to happen")
			return
		}
		if len(code) != 0 {
			//This is a callback from instagram to get the token. We now call a method that retuns us either the token or the error message.
			authToken, er := GetAuthToken(context, w, r, code)
			if er != nil {
				context.Log.Println("Unable to get the token. Error was: ", er)
				return http.StatusInternalServerError, errors.New(er.Error())
			}
			session := context.SessionStore.SessionStart(w, r)
			session.Set("user", authToken.User.FullName)
			session.Set("access_token", authToken.AccessToken)
			//Now redirect the user to the upload page. (ie redirect to the hoemage again, but display the upload template instead)
			return http.StatusSeeOther, nil
		}
		//now that everything's done, we try to render the right template
		// ie upload template if the session has the user token, else the login template
		//lets read the session token
		displayUser := session.Get("user")

		markup := renderIndex(context, displayUser)
		if markup == nil {
			context.Log.Println("Unable to render the templates.")
			return http.StatusInternalServerError, nil
		}
		fmt.Fprint(w, string(markup))
		context.Log.Println("Done generating markup.")
		return http.StatusOK, nil
	default:
		return http.StatusUnauthorized, nil
	}
}

func UploadHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (revVal int, err error) {
	session := context.SessionStore.SessionStart(w, r)
	authToken := session.Get("access_token")
	switch r.Method {
	case "GET":
		// render the template
		context.Log.Println("Rendering the template for upload")
		t, err := template.ParseGlob("templates/upload.html")
		if err != nil {
			context.Log.Println("Unable to parse the Upload template. Error is: ", err)
			return http.StatusInternalServerError, err
		}
		w.Header().Set("Content-Type", "text/html")
		t.Execute(w, nil)
		revVal = http.StatusFound
		err = nil
		return revVal, err
	case "POST":
		if authToken == nil {
			context.Log.Println("Session token is not found the request.")
			revVal = http.StatusUnauthorized
			err = errors.New("You are not allowed to make the request.")
			return
		}
		context.Log.Println("Trying to upload the file now.")

		file, header, err := r.FormFile("file")

		if err != nil {
			context.Log.Println(err)
			return http.StatusInternalServerError, err
		}

		defer file.Close()

		out, err := os.Create("/tmp/uploadedfile")
		if err != nil {
			context.Log.Println("Unable to create the file for writing. Check your write access privilege")
			return http.StatusInternalServerError, err
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			context.Log.Println(w, err)
		}

		context.Log.Println("File uploaded successfully : ")
		context.Log.Println(header.Filename)
		revVal = http.StatusOK
		err = nil
		return revVal, err
	default:
		return http.StatusMethodNotAllowed, errors.New("This method is not allowed on this resource.")
	}
}

func renderIndex(context *AppContext, userWrapper interface{}) []byte {
	// Generate the markup for the index template.
	if userWrapper == nil {
		context.Log.Print("Attempting to render the login template")
		markup := executeTemplate(context, "login", nil)
		if markup == nil {
			return nil
		}
		params := map[string]interface{}{"LayoutContent": template.HTML(string(markup))}
		return executeTemplate(context, "head", params)
	}
	context.Log.Print("Attempting to render the search template")
	markup := executeTemplate(context, "search", nil)
	if markup == nil {
		return nil
	}
	params := map[string]interface{}{"LayoutContent": template.HTML(string(markup))}
	return executeTemplate(context, "head", params)
}
