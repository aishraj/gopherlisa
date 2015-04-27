package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/aishraj/gopherlisa/constants"
)

//TODO: this should be part of the same module as the callback_server. nmae the module authen
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
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// if resp.StatusCode != 200 {
	// 	log.Fatal("Unable to get a 200 from server during code generation for userid: ")
	// }

}
