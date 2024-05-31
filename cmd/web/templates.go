package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.derrc/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet models.Snippet
	Snippets []models.Snippet
	Form any
	Flash string
	IsAuthenticated bool
	CSRFToken string
}

// returns a formatted string representation of time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// store parsed templates in an in-memory cache
func newTemplateCache() (map[string]*template.Template, error) {
	// initialize map to act as cache
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page);

		// register FuncMap with template set before parsing
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// add partials to base template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// add page to template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		
		// cache template set
		cache[name] = ts;
	}

	return cache, nil
}