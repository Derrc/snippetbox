package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// builds servemux for the application with routes/handlers
func (app *application) routes() http.Handler {
	// creates new servemux and registers handler functions for different URL patterns
	mux := http.NewServeMux()

	// creates file server to serve files in 'static' dir
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// middleware chain with our 'standard' middleware used for every request
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}