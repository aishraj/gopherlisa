// +build appengine

package main

import (
	"log"
	"net/http"
)

func init() {
	log.Println("Starting out the server from Google App engine")
	http.HandleFunc("/authorize/", authroizeHandler)
	http.HandleFunc("/", handleBase)
}
