package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/juliflorezg/lets-go/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)

	}

	// files := []string{
	// 	"./ui/html/base.tmpl.html",
	// 	"./ui/html/pages/home.tmpl.html",
	// 	"./ui/html/partials/nav.tmpl.html",
	// }

	// // we get back a template set from file reads
	// // ts, err := template.ParseFiles("./ui/html/pages/home.tmpl.html")
	// ts, err := template.ParseFiles(files...)

	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	// // we know have to use the ExecuteTemplate() method to write the content of the "base" template as the response body.
	// err = ts.ExecuteTemplate(w, "base", nil)

	// if err != nil {
	// 	app.serverError(w, r, err)
	// }
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Method)
	if r.Method != http.MethodPost {
		// We use set to add an "Allow: POST" header to the response header map
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)

		return
	}

	title := "Loguetown Arc"
	content := "After Nami officially joins, the crew heads to the last town before the entrance to the Grand Line, Loguetown,\n the place where Gold Roger was both born and executed. Not only will they have to deal with a powerful Marine\n captain, but also previous enemies looking for revenge."

	expires := 30

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%v", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	fmt.Printf("id: ->%v<-\n", id)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// fmt.Fprintf(w, "Display a specific snippet for ID %d...", id)
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) fooHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Foo"))
}
