package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// new command line flag, name addr, default value :4000
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	// this assigns the value passed on runtime to the addr variable
	// must be used before using the addr variable
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

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

	logger.Info("starting server", "addr", *addr)

	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
