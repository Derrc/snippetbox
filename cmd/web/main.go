package main

import (
	"log"
	"net/http"
)

func main() {
	// creates new servemux and registers handler functions for different URL patterns
	mux := http.NewServeMux()

	// creates file server to serve files in 'static' dir
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("starting server on :4000")

	// listens on the passed TCP network address with our servemux
	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}