// +build appengine

package main

import (
	"log"
	"net/http"

	"github.com/aishraj/gopherlisa/lib/handlers"
)

func init() {
	log.Println("Starting out the server from Google App engine")
	http.HandleFunc("/authorize/", handlers.AuthroizeHandler)
	http.HandleFunc("/", handlers.HandleBase)
}
