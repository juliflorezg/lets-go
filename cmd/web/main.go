package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/juliflorezg/lets-go/internal/models"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	//> this comment style represents the responsibilities of this main function
	// new command line flag, name addr, default value :4000
	//> pass the runtime configuration settings for the application
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "web:web24pass_@@/snippetbox?parseTime=true", "MySQL data source name")

	// this assigns the value passed on runtime to the addr variable
	// must be used before using the addr variable
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db, err := openDB(*dsn)
	defer db.Close()

	if err != nil {
		logger.Error(err.Error())

		os.Exit(1)
	}

	//> Establish the dependencies for the handlers
	app := &application{
		logger: logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("starting server", "addr", *addr)

	//> Run the HTTP server
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	// to verify that everything is set up correctly we need to
	// use the db.Ping() method to create a connection and check for any errors
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
