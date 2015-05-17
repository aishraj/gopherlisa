package lib

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"testing"
)

func TestMatcher(t *testing.T) {
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
	loadedImage, err := LoadImage("/Users/ge3k/go/src/github.com/aishraj/gopherlisa/downloads/kathmandu/10684246_1100050490011860_943310652_n.jpg")
	if err != nil {
		context.Log.Fatal("Cannot load image from disk.", err)
	}
	matchedImage := FindProminentColour(loadedImage)
	nearestImage := findClosestMatch(context, matchedImage, "kathmandu")
	context.Log.Println("neareset one is", nearestImage)

}
