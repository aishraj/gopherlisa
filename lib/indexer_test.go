package lib

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

func TestIndex(t *testing.T) {
	t.Skip("skipping test mode.")
	sessionStore, err := NewSessionManager("gopherId", 3600)
	if err != nil {
		log.Fatal("Unable to start the session store manager.", err)
	}
	Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info.Println("Starting out the go program")
	db, err := sql.Open("mysql", "root:mysql@/gopherlisa")
	if err != nil {
		log.Fatal("Unable to get a connection to MySQL. Error is: ", err)
	}

	context := &AppContext{Info, sessionStore, db} //TODO: add db connection
	directoryName := "kathmandu"                   //TODO: screw it  for now i'm going to remove this test and revoke access later
	n, err := AddImagesToIndex(context, directoryName)
	if err != nil {
		log.Fatal("ERROR!!!!", err)
	}
	log.Println("The number of images we got are", n)
}
