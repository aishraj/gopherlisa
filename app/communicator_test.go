package app

import (
	"github.com/aishraj/gopherlisa/common"
	"log"
	"os"
	"testing"
)

func TestCommunicator(t *testing.T) {
	t.Skip("skipping test mode.") //TODO: This test does not run anymore. Need to add mock db or in memory db connection to the context.
	t.Log("starting test")
	sessionStore, err := common.NewSessionManager("gopherId", 3600)
	if err != nil {
		log.Fatal("Unable to start the session store manager.", err)
	}
	Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info.Println("Starting out the go program")
	context := &common.AppContext{Info, sessionStore, nil}           //TODO: add db connection
	authToken := "32314616.1fb234f.b65b7c681e404f60aa8aaf2dc6f87e34" //TODO: screw it  for now i'm going to remove this test and revoke access later
	searchTerm := "nepal"
	images, err := LoadImages(context, searchTerm, authToken)
	log.Println("We got the following for images: ", images)
}
