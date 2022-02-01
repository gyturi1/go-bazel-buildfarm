package main

import (
	"log"
	"net/http"

	"github.com/gyturi/go-bazel-buildfarm/pkg"
)

const PORT = "8080"

func main() {

	log.Println("starting server, listening on port " + PORT)

	http.HandleFunc("/", pkg.EchoHandler)
	http.ListenAndServe(":"+PORT, nil)
}
