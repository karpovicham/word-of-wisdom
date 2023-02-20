package quotes_book

import (
	"context"
	"encoding/json"
)

// QuotesBook methods
type QuotesBook interface {
	// GetRandomQuote return random quote from the book
	GetRandomQuote(ctx context.Context) (Quote, error)
}

type Quote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

type Book struct {
	Quotes []Quote `json:"quotes"`
}

func (q *Quote) ToJson() []byte {
	b, _ := json.Marshal(q)
	return b
}
