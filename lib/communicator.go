package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func LoadImages(context *AppContext, searchTerm, authToken string) ([]string, error) {
	context.Log.Println("Trying to load images from instagram now.")
	serverURI := "https://api.instagram.com/v1/tags/" + searchTerm + "/media/recent/"

	uri, err := url.Parse(serverURI)
	if err != nil {
		return nil, errors.New("Unable to parse the URL.")
	}
	data := uri.Query()
	data.Set("access_token", authToken)

	uri.RawQuery = data.Encode()

	urlStr := fmt.Sprintf("%v", uri)
	context.Log.Println("The server URI is: ", urlStr)
	return fetchImages(context, urlStr, authToken)
}

func fetchImages(context *AppContext, serverURI, authToken string) ([]string, error) {
	items := make([]string, 0, 1000)
	urlQueue := make([]string, 10, 500)
	firstURL := serverURI
	urlQueue = append(urlQueue, firstURL)

	for len(urlQueue) > 0 || len(items) <= 50 {
		fetchURL, urlQueue := urlQueue[len(urlQueue)-1], urlQueue[:len(urlQueue)-1]
		responseMap, err := fetchServerResponse(context, fetchURL)
		if err != nil {
			errorMessage := "Oops, there was an error geting the server response"
			context.Log.Println(errorMessage)
			return nil, errors.New(errorMessage)
		}
		responseData := responseMap.Data
		for _, responseMeta := range responseData {
			mediaType := responseMeta.MediaType
			if mediaType == "image" {
				thumbNailURL := responseMeta.Images.Thumbnail.URL
				context.Log.Println("*** Parsed the Response for Image URL:", thumbNailURL, "******")
				items = append(items, thumbNailURL)
			}
		}
		nextURL := responseMap.Pagination.NextURL
		urlQueue = append(urlQueue, nextURL)
	}
	return items, nil
}

func fetchServerResponse(context *AppContext, serverURI string) (APIResponse, error) {

	var responseMap APIResponse

	response, err := http.Get(serverURI)
	if err != nil {
		context.Log.Printf("Unable to get the images from instagram.")
		return responseMap, err
	}
	if response.StatusCode != http.StatusOK {
		context.Log.Println("The response was : ", response)
		context.Log.Println("Unable to get a valid response while trying to laod images")
		context.Log.Println("The error received was: ", response.StatusCode)
		return responseMap, errors.New("Can't get images from instagram.")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		context.Log.Println("Unable to parse HTTP response to body")
		return responseMap, errors.New("Did not get a success while posting on instagram")
	}

	json.Unmarshal(body, &responseMap)
	context.Log.Println("The response body is: ", responseMap)
	return responseMap, nil
}
