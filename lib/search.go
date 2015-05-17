package lib

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

func SearchHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (revVal int, err error) {
	if r.Method == "POST" {
		session := context.SessionStore.SessionStart(w, r)
		authToken := session.Get("access_token")
		context.Log.Println("Already authenticated, and therefore fetching the image from instagram now.")
		formData := r.FormValue("searchTerm")
		tokenString, ok := authToken.(string)
		if !ok {
			context.Log.Println("Token cannot be cast to string. ERROR.")
			return http.StatusInternalServerError, errors.New("Token cannot be cast to string")
		}
		context.Log.Println("^^^^^^ The search query is ^^^^^", formData)
		if len(formData) < 1 {
			context.Log.Fatal("Cannot continue, the form data cannot have a length less than 1")
		}
		//small trick, if the directory for the searchterm Already has the required number of files just skip
		if isDownloadRequired(context, formData) == true {
			context.Log.Println("Download is required, downloading now")
			//TODO: if the files are fewer than download and index. (nice to have)
			if !directoryExists(context, formData) {
				err := os.Mkdir("/Users/ge3k/go/src/github.com/aishraj/gopherlisa/downloads/"+formData, 0777)
				if err != nil {
					context.Log.Println("not able to create the directory ***************")
				}
			}

			images, err := LoadImages(context, formData, tokenString)
			if err != nil {
				context.Log.Println("Error fetching from instagram.")
				return http.StatusInternalServerError, err
			}
			context.Log.Println("List of Images we got are:", images)
			downloadCount, ok := DownloadImages(context, images, formData)
			if !ok {
				context.Log.Println("Unable to download images to the path")
				return http.StatusInternalServerError, errors.New("Download failed")
			}
			context.Log.Println("Download count was: ", downloadCount)
			n, ok := ResizeImages(context, formData)
			if !ok {
				context.Log.Println("Unable to resize images")
				return http.StatusInternalServerError, errors.New("Resizing images failed")
			}
			context.Log.Println("Number of images resized was ", n)

			n, err = AddImagesToIndex(context, formData)
			if err != nil {
				context.Log.Println("Unable to add images to index", err)
				return http.StatusInternalServerError, err
			}
			context.Log.Println("Number of images indexed was", n)
		}

		userId := session.Get("userId")
		if userId == nil {
			return http.StatusInternalServerError, errors.New("UserId not there in sesion. ERROR")
		}
		fileId, ok := userId.(string)
		if !ok {
			context.Log.Println("Unable to cast the userid from session storage.")
			return http.StatusInternalServerError, errors.New("Cannot cast the user id from session storage.")
		}

		CreateMosaic(context, fileId, formData)
		//image TODO get resize working first

		//once resized images are there, lets index them.
		return http.StatusOK, nil //TODO: change this to redirect to the mosaic "creating" page (some loading bar or sth)
	}
	return http.StatusMethodNotAllowed, errors.New("This method is not allowed here.")

}

func isDownloadRequired(context *AppContext, searchTerm string) bool {
	files, err := ioutil.ReadDir("/Users/ge3k/go/src/github.com/aishraj/gopherlisa/downloads/" + searchTerm)
	if err != nil {
		context.Log.Println("ERROR: Unable to count the number of files")
		return true //yes download
	}
	fileCount := len(files)
	context.Log.Println("The number of files in the tag dir is", fileCount)
	if fileCount > 1000 { //TODO change it to less than when testing is over
		return true
	}
	return false
}

func directoryExists(context *AppContext, dirname string) bool {
	src, err := os.Stat("/Users/ge3k/go/src/github.com/aishraj/gopherlisa/downloads/" + dirname)
	if err != nil {
		context.Log.Println("Unable to verify OS stat.")
		return false
	}

	// check if the source is indeed a directory or not
	if !src.IsDir() {
		return false
	}
	return true
}
