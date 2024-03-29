package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title, content string, expires int) (int, error)
	Get(id int) (Snippet, error)
	Latest() ([]Snippet, error)
}

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
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (Snippet, error) {
	// return Snippet{}, nil

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?;`

	row := sm.DB.QueryRow(stmt, id)

	var s Snippet

	// Here, we use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. The arguments
	// to row.Scan are *pointers* to the place we want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (sm *SnippetModel) Latest() ([]Snippet, error) {
	// return nil, nil

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP ORDER BY id DESC LIMIT 10`

	rows, err := sm.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* we check for an error from the Query()
	// method. Otherwise, if Query() returns an error, we'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
