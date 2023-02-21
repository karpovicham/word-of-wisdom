package resolver

import (
	"fmt"

	"github.com/karpovicham/word-of-wisdom/internal/proto"
	"github.com/karpovicham/word-of-wisdom/pkg/pow"
)

func (r *resolver) RequestChallenge() (pow.Data, error) {
	if err := r.Msgr.Send(&proto.Message{
		Type: proto.TypeChallenge,
		Data: nil,
	}); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	resp, err := r.Msgr.Receive()
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	if resp.Type != proto.TypeChallenge {
		return nil, ErrInvalidRespType
	}

	if err := checkRespMsgError(resp.Error); err != nil {
		return nil, err
	}

	return resp.Data, nil
}
