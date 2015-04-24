package generator

import (
	"errors"
	"github/aishraj/gopherlisa/models"
	"log"
	"net/http"
)

const serverUrl = ""

//Fetch an array of structured data from Instagram
//Note this method has to be run a go routine
func instaFetch(userName, tag, dataType) (metadata []ImageMetaData) {
	//check in memory if there is an access token for this usenName.
	//If not fetch it
	accessToken, err := getSavedToken(userName)
	if err != nil {
		log.Fatal("Unable to find access token")
	}
	if len(accessToken) == 0 {
		accessToken, err = generateAccessToken(userName)
	}
	//next we make a series of calls to the api until
	// 1. we get the number of images we want
	// 2. the server says we've reached the end of trail (TODO: figure this one out)
	// 3. We bump into an error (any server error, or invalid metadata error)
	// Note here: rather than doing it all in a single thread,
	// sending each request from a go routine until any one reaches
	// the above three condition maybe a way
	// right now doing this in a loop
}

func getSavedToken(userName) (token string, error err) {
	//if its in the db return it, else return empty string
	return "", nil
}

func generateAccessToken(username) {

}
