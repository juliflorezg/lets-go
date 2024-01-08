package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func main() { 
  // Here we use the http.NewServeMux() fun to initialize a new servermux(router), then
  // register the home function as the handler for the "/" URL route/pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)


  log.Println("starting server on port :4000")


  err:= http.ListenAndServe(":4000", mux)
  log.Fatal(err)
}
