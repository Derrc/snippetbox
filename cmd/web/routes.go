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

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// middleware chain with our 'standard' middleware used for every request
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}