package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		// http.NotFound(w, r)
		app.notFound(w)
		// important to return from handler, otherwise it would keep executing
		// and write "Hello from Snippetbox" message to the response
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	// we get back a template set from file reads
	// ts, err := template.ParseFiles("./ui/html/pages/home.tmpl.html")
	ts, err := template.ParseFiles(files...)

	if err != nil {
		// app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		// log.Print(err.Error())
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.serverError(w, r, err)
		return
	}

	// we use execute method to write the template to the response body. Second parameter represents dynamic data
	// err = ts.Execute(w, nil)

	// we know hace to use the ExecuteTemplate() method to write the content of the "base" template as the response body.
	err = ts.ExecuteTemplate(w, "base", nil)

	if err != nil {
		// log.Print(err.Error())
		// app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.serverError(w, r, err)
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Method)
	if r.Method != http.MethodPost {
		// We use set to add an "Allow: POST" header to the response header map
		w.Header().Set("allow", http.MethodPost)

		// w.WriteHeader(405)
		// w.WriteHeader(404) // we can only use WriteHeader once per response and any subsequent try to change the status code once it has changed won't succeed, we get the error 2024/01/09 12:27:45 http: superfluous response.WriteHeader call from main.snippetCreate (main.go:25)
		// w.Header().Set("allow", "POST") //! won't work, must be called before any WriteHeader() or Write()
		// w.Write([]byte("Method not allowed"))

		// this line is a shortcut for WriteHeader & Write above
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)

		return
	}

	// w.Header().Set("Cache-Control", "public, max-age=31536000")
	// w.Header().Add("Cache-Control", "public")
	// w.Header().Add("cache-control", "max-age=31536000")

	w.Write([]byte("Create a new snippet..."))

	// fmt.Println(w.Header())
	// fmt.Println(w.Header().Get("Cache-Control"))
	// fmt.Println(w.Header().Values("Cache-Control"))
	// fmt.Println(len(w.Header().Values("Cache-Control")))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	fmt.Printf("id: ->%v<-\n", id)
	if err != nil || id < 1 {
		// fmt.Println("error:", err)
		// http.Error(w, "Bad Request", http.StatusBadRequest)
		app.notFound(w)
		return
	}
	// Use the fmt.Fprintf() function to interpolate the id value with our response
	// and write it to the http.ResponseWriter.
	fmt.Fprintf(w, "Display a specific snippet for ID %d...", id)
}

func (app *application) fooHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Foo"))
}
