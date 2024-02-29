package mocks

import (
	"time"

	"github.com/juliflorezg/lets-go/internal/models"
)

var mockSnippet = models.Snippet{
	ID:      1,
	Title:   "Sample Snippet 1",
	Content: "Sample content for snippet 1",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (sm *SnippetModel) Insert(title, content string, expires int) (int, error) {
	return 2, nil
}

func (sm *SnippetModel) Get(id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (sm *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet, mockSnippet, mockSnippet}, nil
}
