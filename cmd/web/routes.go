package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/juliflorezg/lets-go/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Here we use the http.NewServeMux() fun to initialize a new servermux(router), then
	// register the app.home method as the handler for the "/" URL route/pattern
	// mux := http.NewServeMux()

	router := httprouter.New()

	// Create a handler function which wraps our notFound() helper, and then
	// assign it as the custom handler for 404 Not Found responses.
	// Can also be defined for 405 Method Not Allowed by using router.MethodNotAllowed
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	//~ Take the ui.Files embedded filesystem and convert it to a http.FS type so
	//~ that it satisfies the http.FileSystem interface. We then pass that to the
	//~ http.FileServer() function to create the file server HTTP handler
	fileServer := http.FileServer(http.FS(ui.Files))
	//~ Our static files are contained in the "static" folder of the ui.Files
	//~ embedded filesystem. So, for example, our CSS stylesheet is located at
	//~ "static/css/main.css". This means that we no longer need to strip the
	//~ prefix from the request URL -- any requests that start with /static/ can
	//~ just be passed directly to the file server and the corresponding static
	//~ file will be served (so long as it exists).
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

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

	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.

	// Use the nosurf middleware on all our 'dynamic' routes.
	dynamicMd := alice.New(app.sessionManager.LoadAndSave, app.noSurf, app.authenticate)

	// Update these routes to use the new dynamic middleware chain followed by
	// the appropriate handler function. Note that because the alice ThenFunc()
	// method returns a http.Handler (rather than a http.HandlerFunc) we also
	// need to switch to registering the route using the router.Handler() method.
	router.Handler(http.MethodGet, "/", dynamicMd.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamicMd.ThenFunc(app.snippetView))

	// routes for user authentication
	router.Handler(http.MethodGet, "/user/signup", dynamicMd.ThenFunc(app.userSignUp))
	router.Handler(http.MethodPost, "/user/signup", dynamicMd.ThenFunc(app.userSignUpPost))
	router.Handler(http.MethodGet, "/user/login", dynamicMd.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamicMd.ThenFunc(app.userLoginPost))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protectedMd := dynamicMd.Append(app.requireAuthentication)

	// Because the 'protected' middleware chain appends to the 'dynamic' chain
	// the noSurf middleware will also be used on the three routes below too.
	router.Handler(http.MethodGet, "/snippet/create", protectedMd.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protectedMd.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protectedMd.ThenFunc(app.userLogoutPost))

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	mdChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Wrap the existing chain with the recoverPanic middleware.
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	// Return the 'standard' middleware chain followed by the router
	return mdChain.Then(router)
}
