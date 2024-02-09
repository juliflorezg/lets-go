package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/juliflorezg/lets-go/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object
func humanDate(t time.Time) string {
	// that date must be used (https://pkg.go.dev/time@go1.21.6#Time.Format)
	return t.Format("02 Jan 2006 at 15:04 MST")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// This map will act as the template cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all filepaths that
	// match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
	// us a slice of all the filepaths for our application 'page' templates
	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name (like 'home.tmpl.html') from the full filepath
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		// The template.FuncMap must be registered with the template set before we
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() *on this template set* to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() *on this template set* to add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts // keys would be like: home.tmpl.html
	}

	return cache, nil
}
