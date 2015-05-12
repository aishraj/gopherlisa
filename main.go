package main

import (
	"database/sql"
	"github.com/aishraj/gopherlisa/lib"
	"log"
	"net/http"
	"os"
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
	db, err := sql.Open("mysql", "root:mysql@/gopherlisa")
	if err != nil {
		log.Fatal("Unable to get a connection to MySQL. Error is: ", err)
	}

	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Images (id INTEGER PRIMARY KEY, image TEXT UNIQUE, used INTEGER, red INTEGER, green INTEGER, blue INTEGER)")
	if err != nil {
		log.Fatal("Unable to create table in db.")
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Colors
        (id INTEGER PRIMARY KEY,
        image_id INTEGER,
        pos INTEGER,
        red INTEGER,
        green INTEGER,
        blue INTEGER)`)
	if err != nil {
		log.Fatal("Unable to create table in db.")
	}
	context := &lib.AppContext{Info, sessionStore, db}
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
