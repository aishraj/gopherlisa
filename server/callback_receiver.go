package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/aishraj/gopherlisa/constants"
)

func StartResponseServer() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	if instaError := r.URL.Query().Get("error"); len(code) != 0 {
		//now we have the code from instagram
		fmt.Fprintf(w, "Hi there, we trying to verfiy you now. Hang on tight.")
		performPostReqeust(code)

	} else if instaError == "access_denied" && r.URL.Query().Get("error_reason") == "user_denied" {
		fmt.Fprintf(w, "Oops, something went wrong.")
	}

}

func performPostReqeust(code string) {
	/*
	  curl -F 'client_id=CLIENT_ID' \
	   -F 'client_secret=CLIENT_SECRET' \
	   -F 'grant_type=authorization_code' \
	   -F 'redirect_uri=AUTHORIZATION_REDIRECT_URI' \
	   -F 'code=CODE' \
	   https://api.instagram.com/oauth/access_token
	*/
	uri, err := url.ParseRequestURI(constants.OauthBaseURI)
	if err != nil {
		log.Fatal("Unable to parse the post uri.")
	}
	uri.Path = "/access_token/"
	data := url.Values{}
	data.Set("client_id", constants.InstagramClientID)
	data.Add("client_secret", constants.InstagramSecret)
	data.Add("grant_type", "authorization_code")
	data.Add("redirect_uri", constants.RedirectURI)

	urlStr := fmt.Sprintf("%v", uri)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal("Unable to send the post request with the code")
	}
	if resp.StatusCode != 200 {
		log.Fatal("Did not get 200 for the post authn request.")
	}
	rawData := make([]byte, 2048) //TODO: check this, just guess it won't be more than 2 kb.
	n, err := resp.Body.Read(rawData)
	if err != nil {
		log.Fatal("Unable to read from the response buffer of the post data.")
	}
	rawData = rawData[:n] //TODO: check this as well.
}
