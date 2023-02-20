package resolver

import (
	"github.com/karpovicham/word-of-wisdom/internal/proto"
)

func (r *apiResolver) Stop() error {
	return r.Msgr.Send(&proto.Message{
		Type: proto.Stop,
		Data: nil,
	})
}
