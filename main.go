package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aishraj/gopherlisa/lib"
)

func init() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func main() {
	sessionStore, err := lib.NewSessionManager("gopherId", 3600)
	if err != nil {
		log.Fatal("Unable to start the session store manager.", err)
	}

	Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info.Println("Starting out the go program")
	context := &lib.AppContext{Info, sessionStore}
	authHandler := lib.Handler{context, lib.AuthroizeHandler}
	rootHandler := lib.Handler{context, lib.BaseHandler}
	uploadHandler := lib.Handler{context, lib.UploadHandler}
	searchHandler := lib.Handler{context, lib.SearchHandler}
	http.Handle("/login/", authHandler)
	http.Handle("/search", searchHandler)
	http.Handle("/upload/", uploadHandler)
	http.Handle("/", rootHandler)
	http.ListenAndServe(":8080", nil)
}
