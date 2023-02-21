package resolver

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/karpovicham/word-of-wisdom/internal/proto"
	"github.com/karpovicham/word-of-wisdom/pkg/pow"
	"github.com/karpovicham/word-of-wisdom/service/quotes_book"
)

var (
	ErrInvalidRespData  = errors.New("invalid response data")
	ErrInvalidRespType  = errors.New("invalid response type")
	ErrUnknownRespError = errors.New("invalid response error")
	ErrInvalidReqData   = errors.New("invalid request data")
	ErrNotVerified      = errors.New("not verified")
)

func (r *resolver) RequestQuote(powData pow.Data) (*quotes_book.Quote, error) {
	if err := r.Msgr.Send(&proto.Message{
		Type: proto.TypeQuote,
		Data: powData,
	}); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	resp, err := r.Msgr.Receive()
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	if resp.Type != proto.TypeQuote {
		return nil, ErrInvalidRespType
	}

	if err := checkRespMsgError(resp.Error); err != nil {
		return nil, err
	}

	var quote quotes_book.Quote
	if err := json.Unmarshal(resp.Data, &quote); err != nil {
		return nil, ErrInvalidRespData
	}

	return &quote, nil
}
