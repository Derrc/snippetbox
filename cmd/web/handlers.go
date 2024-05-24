package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// handler function
func home(w http.ResponseWriter, r *http.Request) {
	// adds a header 'Server: Go' to the response header map
	w.Header().Add("Server", "Go")

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// execute 'base' template
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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