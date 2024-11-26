package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
)

const portNumOne = ":2222"
const portNumTwo = ":4444"

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
	log.Println("Listening on ports", portNumOne, portNumTwo)

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/about", about)

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
