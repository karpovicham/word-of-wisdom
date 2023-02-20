package resolver

import (
	"fmt"

	"github.com/karpovicham/word-of-wisdom/internal/proto"
	"github.com/karpovicham/word-of-wisdom/pkg/pow"
)

func (r *apiResolver) RequestChallenge() (pow.Data, error) {
	if err := r.Msgr.Send(&proto.Message{
		Type: proto.Challenge,
		Data: nil,
	}); err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}

	resp, err := r.Msgr.Receive()
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	if resp.Type != proto.Challenge {
		return nil, ErrInvalidRespType
	}

	return resp.Data, nil
}
