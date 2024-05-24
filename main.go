package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// handler function
func home(w http.ResponseWriter, r *http.Request) {
	// adds a header 'Server: Go' to the response header map
	w.Header().Add("Server", "Go")

	w.Write([]byte("Hello from Snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	// extract id wildcard from request url and try to convert to integer
	// return 404 page not found if wildcard value is not an integer or less than 1
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// interpolate string and write to passed Writer (ResponseWriter implements Writer)
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// sends a 201 status code with the response
	w.WriteHeader(http.StatusCreated);

	w.Write([]byte("Save a new snippet..."))
}

func main() {
	// creates new servemux and registers handler functions for different URL patterns
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on :4000")

	// listens on the passed TCP network address with our servemux
	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}