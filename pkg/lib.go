package pkg

import (
	"fmt"
	"log"
	"net/http"
)

func Format(e error) string {
	return fmt.Sprintf("couldn't accept: %s", e)
}

// EchoHandler echos back the request as a response
func EchoHandler(writer http.ResponseWriter, request *http.Request) {

	log.Println("Echoing back request made to " + request.URL.Path + " to client (" + request.RemoteAddr + ")")

	request.Write(writer)
}
