// +build appengine

package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aishraj/gopherlisa/lib"
	_ "github.com/aishraj/gopherlisa/lib/memory"
	"github.com/aishraj/gopherlisa/lib/sessions"
)

func init() {
	sessionStore, err := sessions.NewManager("memory", "gosessionid", 3600)
	if err != nil {
		log.Fatal("Unable to start the session store manager.", err)
	}

	Info := log.New(ioutil.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info.Println("Starting out the go program")
	context := &lib.AppContext{Info, sessionStore}
	authHandler := lib.Handler{context, lib.AuthroizeHandler}
	rootHandler := lib.Handler{context, lib.HandleBase}
	http.Handle("/authorize/", authHandler)
	http.Handle("/", rootHandler)
}
