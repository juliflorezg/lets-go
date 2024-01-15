package models

import (
	"database/sql"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet.
// The fields of the struct correspond to the fields in the MySQL snippets table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
  DB *sql.DB
}

func (sm *SnippetModel) Insert(title, content string, expires int) (int, error) {

	// the SQL statement we want to execute on the DB
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP, INTERVAL ? DAY))`

	result, err := sm.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil{
		return 0, err
	}

	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (Snippet, error){
  return Snippet{}, nil
}

func (sm *SnippetModel) Latest() ([]Snippet, error) {
  return nil, nil
}
