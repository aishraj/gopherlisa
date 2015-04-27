package main

import (
	"log"
	"net/http"

	"github.com/aishraj/gopherlisa/handlers"
)

func handleBase(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside the base handler.")
	handlers.BaseHandler(w, r)
}

func authroizeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inside the authorization handler.")
	handlers.AuthHandler(w, r)
}

func main() {
	http.HandleFunc("/authorize/", authroizeHandler)
	http.HandleFunc("/", handleBase)
	http.ListenAndServe(":8000", nil)
}
