// +build !appengine

package main

import (
	"net/http"

	"github.com/aishraj/gopherlisa/handlers"
)

func main() {
	http.HandleFunc("/authorize/", handlers.AuthroizeHandler)
	http.HandleFunc("/", handlers.HandleBase)
	http.ListenAndServe(":8000", nil)
}
