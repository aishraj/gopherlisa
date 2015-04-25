package generator

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/aishraj/gopherlisa/constants"
)

func AuthenticateUser(user string) {
	u, err := url.Parse(constants.OauthBaseURI)
	u.Path = "/authorize/"
	if err != nil {
		log.Fatal("Unable to parse the URI. Error is: ", err)
	}
	query := u.Query()
	query.Set("client_id", constants.InstagramClientID)
	query.Set("redirect_uri", constants.RedirectURI)
	query.Set("response_type", "code")

	u.RawQuery = query.Encode()

	resp, err := http.Get(fmt.Sprintf("%v", u)) //TODO check this
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("Unable to get a 200 from server during code generation for userid: ", user)
	}

}
