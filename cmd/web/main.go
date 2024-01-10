package main

import (
	"log"
	"net/http"
)

func main() {
	// Here we use the http.NewServeMux() fun to initialize a new servermux(router), then
	// register the home function as the handler for the "/" URL route/pattern
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// we use the the created fileServer as a handler for any request coming to /static/
	// for matching paths, we strip (remove) the '/static' from the path before it reaches the fileServer handler so if can give back the correct file.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	// Here we're registering two new handler functions and corresponding URL patterns with the servermuxx
	mux.HandleFunc("/snippet/create", snippetCreate)
	mux.HandleFunc("/snippet/view", snippetView)
	//* by including a host name in the URL pattern, we can route requests based on the host part of the URL,
	//* the handler would be invoked for requests like 'http://snippet.view.org/anypath'
	// mux.HandleFunc("snippet.view.org/", snippetView)
	//subtree path, if we make a request to /foo it will automatically redirect to /foo/
	mux.HandleFunc("/foo/", fooHandler)

	log.Println("starting server on port :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
