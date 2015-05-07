package lib

import (
	"net/http"
)

func searchHandler(context *AppContext, w http.ResponseWriter, r *http.Request) (revVal int, err error) {
	return http.StatusTeapot, nil
}
