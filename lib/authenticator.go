// +build !appengine

package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/aishraj/gopherlisa/lib/constants"
)

func AuthenticateUser(context *AppContext, w http.ResponseWriter, r *http.Request) (status int, err error) {
	method := r.Method

	switch method {

	case "POST":
		u, err := url.Parse(constants.OauthBaseURI)
		u.Path = "/oauth/authorize/"
		if err != nil {
			context.Log.Println("Unable to parse the URI. Error is: ", err)
			return http.StatusBadRequest, err
		}
		query := u.Query()
		query.Set("client_id", constants.InstagramClientID)
		query.Set("redirect_uri", constants.RedirectURI)
		query.Set("response_type", "code")

		u.RawQuery = query.Encode()
		context.Log.Println("Query is: ", u)
		http.Redirect(w, r, fmt.Sprintf("%v", u), http.StatusSeeOther) //TODO change this so that redirect happens in the calling method.
		return http.StatusOK, nil                                      //TODO change this
	default:
		context.Log.Print("This method is not allowed")
		return http.StatusMethodNotAllowed, errors.New("Not allowed.")
	}
}

func PerformPostReqeust(applicationContext *AppContext, w http.ResponseWriter, req *http.Request, code string) (status int, err error) {

	applicationContext.Log.Printf("Performing Post trigggered with the code value %v \n", code)

	uri, err := url.ParseRequestURI(constants.OauthBaseURI)

	if err != nil {
		return http.StatusInternalServerError, err
	}

	uri.Path = "/oauth/access_token/"
	data := url.Values{}
	data.Set("client_id", constants.InstagramClientID)
	data.Add("client_secret", constants.InstagramSecret)
	data.Add("grant_type", "authorization_code")
	data.Add("redirect_uri", constants.RedirectURI)
	data.Add("code", code)

	urlStr := fmt.Sprintf("%v", uri)

	applicationContext.Log.Print("posting to the url: ", urlStr)

	client := &http.Client{}

	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	resp, err := client.Do(r)

	if err != nil {
		applicationContext.Log.Println("Unable to send the post request with the code")
		return http.StatusInternalServerError, err

	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		applicationContext.Log.Println("Did not get 200 for the post authn request.")
		return resp.StatusCode, errors.New("Did not get a success while posting on instagram")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		applicationContext.Log.Println("Unable to parse HTTP response to body")
		return http.StatusInternalServerError, errors.New("Did not get a success while posting on instagram")
	}

	var authToken AuthenticationResponse

	err = json.Unmarshal(body, &authToken)
	if err != nil {
		applicationContext.Log.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			applicationContext.Log.Println(string(body[v.Offset-40 : v.Offset]))
		}
		applicationContext.Log.Println("Eror while unmarshalling data.")
		return http.StatusInternalServerError, errors.New("Error unmarshalling data")
	}

	applicationContext.Log.Printf("Yippie!! your authentication token is %v \n", authToken.AccessToken)
	applicationContext.Log.Println("Got the data: ", authToken)

	session := applicationContext.SessionStore.SessionStart(w, r)
	session.Set("user", authToken.User.FullName)
	session.Set("access_token", authToken.AccessToken)
	// TODO: set data in sesison storeage.
	//key, err := datastore.Put(context, datastore.NewIncompleteKey(context, "authToken", nil), &authToken)
	// applicationContext.log.Printf("Key is %v \n", key)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	return http.StatusSeeOther, nil //TODO should work. lets see what happens

	//TODO save this data along with the user info info db. Override the token if user is already present.

	//Now redirect to homepage (where we again check the token and this time it will be around)
}
