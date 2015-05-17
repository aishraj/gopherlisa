package main

import (
	"database/sql"
	"github.com/aishraj/gopherlisa/lib"
	_ "github.com/go-sql-driver/mysql"
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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Images ( id int(5) NOT NULL AUTO_INCREMENT, imgtype varchar(255),  img varchar(255) UNIQUE, red int(16), green int(16), blue int(16), PRIMARY KEY(id) )")
	if err != nil {
		log.Fatal("Unable to create table in db. Aborting now. Error is :", err)
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
