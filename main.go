package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		// important to return from handler, otherwise it would keep executing
		// and write "Hello from Snippetbox" message to the response
		return
	}
	w.Write([]byte("Hello from Snippetbox"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Method)
	if r.Method != "POST" {
		// We use set to add an "Allow: POST" header to the response header map
		w.Header().Set("allow", "POST") 
		
		// w.WriteHeader(405)
		// w.WriteHeader(404) // we can only use WriteHeader once per response and any subsequent try to change the status code once it has changed won't succeed, we get the error 2024/01/09 12:27:45 http: superfluous response.WriteHeader call from main.snippetCreate (main.go:25)
		// w.Header().Set("allow", "POST") // ! won't work, must be called before any WriteHeader() or Write()
		// w.Write([]byte("Method not allowed")) 

		// this line is a shortcut for WriteHeader & Write above
		http.Error(w, "Method not allowed", 405)
		
		return
	}

	w.Write([]byte("Create a new snippet..."))
}
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}
func fooHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Foo"))
}

func main() {
	// Here we use the http.NewServeMux() fun to initialize a new servermux(router), then
	// register the home function as the handler for the "/" URL route/pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	// Here we're registering two new handler functions and corresponding URL patterns with the servermux
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

	// v:= bird object
}
