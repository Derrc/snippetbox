package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.derrc/internal/models"
	"snippetbox.derrc/ui"
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

	// get all matching filepaths from embedded filesystem
	pages, err := fs.Glob(ui.Files, "./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page);

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		// parse template files that match patterns
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// cache template set
		cache[name] = ts;
	}

	return cache, nil
}