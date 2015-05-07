package lib

import (
	"errors"
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
func baseHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (retVal int, err error) {
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
			} else {
				session := context.SessionStore.SessionStart(w, r)
				session.Set("user", authToken.User.FullName)
				session.Set("access_token", authToken.AccessToken)
			}
			//Now redirect the user to the upload page. (ie redirect to the hoemage again, but display the upload template instead)
			return http.StatusSeeOther, nil
		}
		//now that everything's done, we try to render the right template
		// ie upload template if the session has the user token, else the login template
	}
}

func oldBaseHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (retVal int, err error) {
	context.Log.Println("Inside the base handler.")
	session := context.SessionStore.SessionStart(w, r)
	context.Log.Println("Got session object inside the base handler.")
	createtime := session.Get("createtime")
	context.Log.Println("The create time of the session is: ", createtime)

	if createtime == nil {
		session.Set("createtime", time.Now().Unix())
	} else if (createtime.(int64) + 360) < (time.Now().Unix()) {
		context.Log.Println("Session has expired, starting new session.")
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
			//return GetAuthToken(context, w, r, code) //performs a post request to instagram to get the auth token. Should redirect upon completion.
			//TODO: This method need not be the return value and can be a generic method.
		} else if instaError == "access_denied" && r.URL.Query().Get("error_reason") == "user_denied" {
			context.Log.Println("User denied permission for access.")
			retVal = http.StatusUnauthorized
			err = errors.New("Seems you didn't allow that to happen")
		} else {
			//This is where we start.
			context.Log.Println("Looks like a new user!. Trying to get the value of the user name from the session")
			displayUser := session.Get("user")
			if displayUser == nil {
				context.Log.Println("The user's name was not set. Setting it to guest now")
				session.Set("user", "Guest")
				//				renderLoginIndex
				//Show template login
			} else {
				context.Log.Println("The user's name was set, displaying it now.")
				session.Set("user", displayUser)
				//Show template search
			}
			context.Log.Println("Now trying to render the template.")
			t, err := template.ParseGlob("templates/login.html")
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
		retVal = http.StatusMethodNotAllowed
		err = errors.New("Method not supported.")
	}
	context.Log.Printf("Returning with values Status  %v and error %v ", retVal, err)
	return
}

func uploadHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (revVal int, err error) {
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

func authroizeHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (revVal int, err error) {
	session := context.SessionStore.SessionStart(w, r)
	authToken := session.Get("access_token")
	if authToken == nil {
		context.Log.Println("Inside the authorization handler.")
		return AuthenticateUser(context, w, r)
	}
	context.Log.Println("Already authenticated, and therefore fetching the image from instagram now.")
	formData := r.FormValue("searchTerm")
	tokenString, ok := authToken.(string)
	if !ok {
		context.Log.Println("Token cannot be cast to string. ERROR.")
		return http.StatusInternalServerError, errors.New("Token cannot be cast to string")
	}
	images, err := LoadImages(context, formData, tokenString)
	if err != nil {
		context.Log.Println("Error fetching from instagram.")
		return http.StatusInternalServerError, err
	}
	context.Log.Println("List of Images we got are:", images)
	downloadCount, ok := DownloadImages(images)
	if !ok {
		context.Log.Println("Unable to download images to the path")
		return http.StatusInternalServerError, errors.New("Download failed")
	}
	context.Log.Println("Download count was: ", downloadCount)
	//imageIndex := buildImageIndex(downloadPath) //traverse this os path and build an index type of index is yet to be decided, but will be most likely db backed.
	//TODO: change the original workflow to allow image upload.
	//FLow -> Welcome guest, sign in to get started
	// Post sign in -> Upload Image (Step 1)
	// Redirect to step 2 -> Search for item
	// Now proceed.
	// Divide the actual image in a grid.
	// Fetch the image which is closest to the average color of the grid.

	return http.StatusSeeOther, nil

}
