package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/aishraj/gopherlisa/constants"
)

func performPostReqeust(code string) {

	uri, err := url.ParseRequestURI(constants.OauthBaseURI)

	if err != nil {
		log.Fatal("Unable to parse the post uri.")
	}

	uri.Path = "/oauth/access_token/"
	data := url.Values{}
	data.Set("client_id", constants.InstagramClientID)
	data.Add("client_secret", constants.InstagramSecret)
	data.Add("grant_type", "authorization_code")
	data.Add("redirect_uri", constants.RedirectURI)
	data.Add("code", code)

	urlStr := fmt.Sprintf("%v", uri)

	log.Print("posting to the url: ", urlStr)

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
