package in_memory_book

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/karpovicham/word-of-wisdom/service/quotes_book"
)

var ErrNoQuotes = errors.New("no quotes")

// NewQuotesBookService returns quotesBookService with loaded books and seeded randomizer
func NewQuotesBookService(quotesFilepath string, randoSource rand.Source) (quotes_book.QuotesBook, error) {
	var book quotes_book.Book
	b, err := os.ReadFile(quotesFilepath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	err = json.Unmarshal(b, &book)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}

	if len(book.Quotes) == 0 {
		return nil, ErrNoQuotes
	}

	return &quotesBookService{
		book:       book,
		randomizer: rand.New(randoSource),
	}, nil
}

type quotesBookService struct {
	// In memory stored quotes collection
	book quotes_book.Book
	// Randomized should be seeded to truly return random values each runtime
	randomizer *rand.Rand
}

// GetRandomQuote returns GetRandomQuote
func (r *quotesBookService) GetRandomQuote(_ context.Context) (quotes_book.Quote, error) {
	randQuoteIndex := r.randomizer.Intn(len(r.book.Quotes))
	return r.book.Quotes[randQuoteIndex], nil
}
