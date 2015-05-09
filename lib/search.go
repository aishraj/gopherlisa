package lib

import (
	"errors"
	"net/http"
)

func SearchHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (revVal int, err error) {
	session := context.SessionStore.SessionStart(w, r)
	authToken := session.Get("access_token")
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
	return http.StatusSeeOther, nil //TODO: change this to redirect to the mosaic "creating" page (some loading bar or sth)
}
