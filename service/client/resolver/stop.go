package resolver

import (
	"github.com/karpovicham/word-of-wisdom/internal/proto"
)

func (r *resolver) Stop() error {
	return r.Msgr.Send(&proto.Message{
		Type: proto.TypeStop,
		Data: nil,
	})
}
