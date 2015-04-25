package generator

import (
	"fmt"
	"log"

	"github.com/aishraj/gopherlisa/models"
)

const serverURL = ""

//Fetch an array of structured data from Instagram
//Note this method has to be run a go routine
func instaFetch(userName string, tag string, dataType string) (metadata []models.ImageMetaData) {
	//check in the database if there is an access token for this usenName.
	//If not fetch it
	accessToken, err := getSavedToken(userName)
	if err != nil {
		log.Fatal("Unable to find access token")
	}
	if len(accessToken) == 0 {
		accessToken, err := generateAccessToken(userName)
		if err != nil {
			log.Fatal("Unable to generate access token. Error is: ", err)
		}
		fmt.Println("access token is: ", accessToken)
	}
	//next we make a series of calls to the api until
	// 1. we get the number of images we want
	// 2. the server says we've reached the end of trail (TODO: figure this one out)
	// 3. We bump into an error (any server error, or invalid metadata error)
	// Note here: rather than doing it all in a single thread,
	// sending each request from a go routine until any one reaches
	// the above three condition maybe a way
	// right now doing this in a loop
	return nil
}

func getSavedToken(userName string) (token string, err error) {
	//if its in the db return it, else return empty string
	return "", nil
}

func generateAccessToken(username string) (accessToken string, err error) {
	//This is where we communicate to instagram.
	//Steps based on https://instagram.com/developer/authentication/

	return "", nil
}
