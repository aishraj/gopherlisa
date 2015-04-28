package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/aishraj/gopherlisa/constants"
	"github.com/aishraj/gopherlisa/models"
)

func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(constants.OauthBaseURI)
	u.Path = "/oauth/authorize/"
	if err != nil {
		log.Fatal("Unable to parse the URI. Error is: ", err)
	}
	query := u.Query()
	query.Set("client_id", constants.InstagramClientID)
	query.Set("redirect_uri", constants.RedirectURI)
	query.Set("response_type", "code")

	u.RawQuery = query.Encode()
	fmt.Println(u)
	http.Redirect(w, r, fmt.Sprintf("%v", u), 301) //TODO check this
}

func PerformPostReqeust(w http.ResponseWriter, req *http.Request, code string) {

	log.Printf("Perform Post trigggered with the code value %v \n", code)

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

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Did not get 200 for the post authn request.")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Unable to parse HTTP response to body")
	}

	var authToken models.AuthenticationResponse

	err = json.Unmarshal(body, &authToken)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println(string(body[v.Offset-40 : v.Offset]))
		}
		log.Fatal("Error unmarshalling data")
	}

	fmt.Printf("Yippie!! your authentication token is %v \n", authToken.AccessToken)

	http.Redirect(w, req, "/", 301) //TODO should work. lets see what happens

	//TODO save this data along with the user info info db. Override the token if user is already present.

	//Now redirect to homepage (where we again check the token and this time it will be around)
}
