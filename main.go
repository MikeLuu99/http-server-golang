package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
)

// Quote structure for JSON
type Quote struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

// Global variables
var (
	quotes     []Quote
	quotesFile = "quotes.json"
	mutex      sync.Mutex
)

const portNumOne = ":2222"
const portNumTwo = ":4444"

type serverAddressType string

const serverAddress serverAddressType = "localhost"

// loadQuotes reads quotes from the JSON file
func loadQuotes() error {
	data, err := os.ReadFile(quotesFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &quotes)
}

// saveQuotes writes quotes to the JSON file
func saveQuotes() error {
	data, err := json.MarshalIndent(quotes, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(quotesFile, data, 0644)
}

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

	log.Printf("%s got / request first(%t):%s second(%t):%s \nbody:%s",
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

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	myName := r.PostFormValue("myName")
	if myName == "" {
		w.Header().Set("x-missing-name", "true")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("x-missing-name", "false")
	log.Printf("%s got /about request \n Hi %s", ctx.Value(serverAddress), myName)
	fmt.Fprintf(w, "About Page")
}

func quotesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		mutex.Lock()
		defer mutex.Unlock()

		if len(quotes) == 0 {
			http.Error(w, "No quotes found", http.StatusNotFound)
			log.Printf("GET /quotes")
			return
		}

		if r.URL.Query().Has("author") {
			author := r.URL.Query().Get("author")
			for _, quote := range quotes {
				if quote.Author == author {
					json.NewEncoder(w).Encode(quote)
					log.Printf("GET /quotes?author=%s", author)
					return
				}
			}
			http.Error(w, "No quotes found for this author", http.StatusNotFound)
			log.Printf("GET /quotes?author=%s", author)
			return
		}

		// If no author specified, return a random quote
		randomIndex := rand.Intn(len(quotes))
		randomQuote := quotes[randomIndex] // You might want to implement random selection here
		json.NewEncoder(w).Encode(randomQuote)
		log.Printf("GET /quotes: random quote: %s", randomQuote)

	case http.MethodPost:
		var newQuote Quote

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			log.Printf("Error reading request body: %v", err)
			return
		}

		err = json.Unmarshal(body, &newQuote)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			log.Printf("Error parsing JSON: %v", err)
			return
		}

		if newQuote.Author == "" || newQuote.Quote == "" {
			http.Error(w, "Author and quote fields are required", http.StatusBadRequest)
			log.Printf("Error parsing JSON: %v", err)
			return
		}

		mutex.Lock()
		quotes = append(quotes, newQuote)
		log.Printf("POST /quotes: added quote: %s", newQuote)

		err = saveQuotes()
		mutex.Unlock()

		if err != nil {
			http.Error(w, "Error saving quote to file", http.StatusInternalServerError)
			log.Printf("Error saving quote to file: %v", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Quote added successfully"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := loadQuotes()
	if err != nil {
		log.Printf("Error loading quotes: %v", err)
		quotes = []Quote{} // Initialize empty if file can't be loaded
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/about", about)
	mux.HandleFunc("/quotes", quotesHandler)

	serverOne := &http.Server{
		Addr:    portNumOne,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return context.WithValue(ctx, serverAddress, l.Addr().String())
		},
	}

	serverTwo := &http.Server{
		Addr:    portNumTwo,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return context.WithValue(ctx, serverAddress, l.Addr().String())
		},
	}

	// Use WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Printf("Server starting on port %s", portNumOne)
		if err := serverOne.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server one error: %v", err)
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		log.Printf("Server starting on port %s", portNumTwo)
		if err := serverTwo.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server two error: %v", err)
			cancel()
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Shutdown both servers gracefully
	if err := serverOne.Shutdown(context.Background()); err != nil {
		log.Printf("Server One shutdown error: %v", err)
	}
	if err := serverTwo.Shutdown(context.Background()); err != nil {
		log.Printf("Server Two shutdown error: %v", err)
	}

	wg.Wait()
}
