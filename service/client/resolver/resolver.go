package resolver

import (
	"github.com/karpovicham/word-of-wisdom/internal/messenger"
	"github.com/karpovicham/word-of-wisdom/internal/proto"
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

type resolver struct {
	Msgr messenger.Messenger
}

// NewClientAPIResolver returns implementer Resolver
func NewClientAPIResolver(msgr messenger.Messenger) Resolver {
	return &resolver{
		Msgr: msgr,
	}
}

// checkRespMsgError return domain like error message proto contains one
func checkRespMsgError(err *proto.Error) error {
	if err == nil {
		return nil
	}

	switch *err {
	case proto.ErrorInvalidData:
		return ErrInvalidReqData
	case proto.ErrorNotVerified:
		return ErrNotVerified
	default:
		return ErrUnknownRespError
	}
}
