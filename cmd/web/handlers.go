package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/juliflorezg/lets-go/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// removing manual check of "/" path since httprouter matches the "/" path exactly
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl.html", templateData)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// checking this method is no longer necessary because the httprouter does this for us
	// fmt.Println(r.Method)
	// if r.Method != http.MethodPost {
	// 	// We use set to add an "Allow: POST" header to the response header map
	// 	w.Header().Set("Allow", http.MethodPost)

	// 	app.clientError(w, http.StatusMethodNotAllowed)

	// 	return
	// }

	title := "Loguetown Arc"
	content := "After Nami officially joins, the crew heads to the last town before the entrance to the Grand Line, Loguetown, the place where Gold Roger was both born and executed. Not only will they have to deal with a powerful Marine captain, but also previous enemies looking for revenge."

	expires := 30

	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	// Update the redirect path to use the new clean URL format.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
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

	templateData := app.newTemplateData(r)
	templateData.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl.html", templateData)
}

func (app *application) fooHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Foo"))
}
