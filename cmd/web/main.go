package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	//> this comment style represents the responsibilities of this main function
	// new command line flag, name addr, default value :4000
	//> pass the runtime configuration settings for the application
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	// this assigns the value passed on runtime to the addr variable
	// must be used before using the addr variable
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	//> Establish the dependencies for the handlers
	app := &application{
		logger: logger,
	}

	logger.Info("starting server", "addr", *addr)

	//> Run the HTTP server
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
