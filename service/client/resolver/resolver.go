package resolver

import (
	"github.com/karpovicham/word-of-wisdom/internal/messenger"
	"github.com/karpovicham/word-of-wisdom/pkg/pow"
	"github.com/karpovicham/word-of-wisdom/service/quotes_book"
)

// Resolver represents basic API Resolver for Client-Server requests
type Resolver interface {
	// RequestChallenge request Challenge, receive the Data for the POW work to be done
	RequestChallenge() (pow.Data, error)
	// RequestQuote request quote with the work done Data
	RequestQuote(powData pow.Data) (*quotes_book.Quote, error)
	// Stop sends signal to client to close connection
	Stop() error
}

type apiResolver struct {
	Msgr messenger.Messenger
}

func NewClientAPIResolver(msgr messenger.Messenger) Resolver {
	return &apiResolver{
		Msgr: msgr,
	}
}
