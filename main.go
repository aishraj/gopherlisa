package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aishraj/gopherlisa/lib"
	_ "github.com/aishraj/gopherlisa/lib/memory"
	"github.com/aishraj/gopherlisa/lib/sessions"
)

func main() {
	sessionStore, err := sessions.NewManager("memory", "gosessionid", 3600)
	if err != nil {
		log.Fatal("Unable to start the session store manager.", err)
	}

	Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info.Println("Starting out the go program")
	context := &lib.AppContext{Info, sessionStore}
	authHandler := lib.Handler{context, lib.AuthroizeHandler}
	rootHandler := lib.Handler{context, lib.HandleBase}
	http.Handle("/authorize/", authHandler)
	http.Handle("/", rootHandler)
	http.ListenAndServe(":8080", nil)
}
