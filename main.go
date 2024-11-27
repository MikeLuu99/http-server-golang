package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

const portNumOne = ":2222"
const portNumTwo = ":4444"

type serverAddressType string

const serverAddress serverAddressType = "localhost"

func home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%sgot / request first(%t):%s second(%t):%s \nbody:%s",
		ctx.Value(serverAddress),
		hasFirst,
		first,
		hasSecond,
		second,
		body)
	fmt.Fprintf(w, "Home Page")
}

func about(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	myName := r.PostFormValue("myName")
	if myName == "" {
		w.Header().Set("x-missing-name", "true")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		w.Header().Set("x-missing-name", "false")
		w.WriteHeader(http.StatusOK)
		log.Printf("%s got /about request \n Hi %s", ctx.Value(serverAddress), myName)
	}
	fmt.Fprintf(w, "About Page")
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "text/plain")

	response := fmt.Sprintf("Request Method: %s\nMessage: Hello from the server!", r.Method)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func main() {
	log.Println("Listening on ports", portNumOne, portNumTwo)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/about", about)
	mux.HandleFunc("/get", handleGet)

	ctx := context.Background()
	serverOne := &http.Server{
		Addr:    portNumOne,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, serverAddress, l.Addr().String())
			return ctx
		},
	}

	serverTwo := &http.Server{
		Addr:    portNumTwo,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, serverAddress, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server closed under request")
		} else if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err := serverTwo.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server closed under request")
		} else if err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
