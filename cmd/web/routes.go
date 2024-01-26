package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Here we use the http.NewServeMux() fun to initialize a new servermux(router), then
	// register the app.home method as the handler for the "/" URL route/pattern
	// mux := http.NewServeMux()

	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// we use the created fileServer as a handler for any request coming to /static/
	// for matching paths, we strip (remove) the '/static' from the path before it reaches the fileServer handler so if can give back the correct file.
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// mux.HandleFunc("/", app.home)
	// Here we're registering two new handler functions and corresponding URL patterns with the servermux
	// mux.HandleFunc("/snippet/create", app.snippetCreate)
	// mux.HandleFunc("/snippet/view", app.snippetView)
	//* by including a host name in the URL pattern, we can route requests based on the host part of the URL,
	//* the handler would be invoked for requests like 'http://snippet.view.org/anypath'
	// mux.HandleFunc("snippet.view.org/", snippetView)
	//subtree path, if we make a request to /foo it will automatically redirect to /foo/
	// mux.HandleFunc("/foo/", app.fooHandler)

	// Create a handler function which wraps our notFound() helper, and then
	// assign it as the custom handler for 404 Not Found responses.
	// Can also be defined for 405 Method Not Allowed by using router.MethodNotAllowed
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	mdChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Wrap the existing chain with the recoverPanic middleware.
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// Return the 'standard' middleware chain followed by the router
	return mdChain.Then(router)
}
