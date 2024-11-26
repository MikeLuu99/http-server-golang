package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
)

const serverAddress = "localhost"

func home(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	log.Printf("%sgot / request first(%t):%s second(%t):%s ", ctx.Value(serverAddress), hasFirst, first, hasSecond, second)
	fmt.Fprintf(w, "Home Page")
}

func about(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("request to about", ctx.Value(serverAddress))
	fmt.Fprintf(w, "About Page")
}

func main() {
	log.Println("Listening on port 3333")

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/about", about)

	ctx := context.Background()
	serverOne := &http.Server{
		Addr:    ":3333",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, serverAddress, l.Addr().String())
			return ctx
		},
	}

	err := serverOne.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed under request")
	} else if err != nil {
		log.Fatal(err)
	}
}
