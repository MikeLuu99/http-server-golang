package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNum = ":8080"

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About Page")
}

func main() {
	log.Println("Listening on port", portNum)

	http.HandleFunc("/", home)
	http.HandleFunc("/about", about)

	log.Println("Starting server on port", portNum)

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}
