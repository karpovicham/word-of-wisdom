package hashcash

import (
	"context"
	"errors"
	"fmt"
	"hash"

	"github.com/karpovicham/word-of-wisdom/pkg/pow"
)

var (
	// ErrInvalidWorkDone returned when worked checks result work data and it's not valid
	ErrInvalidWorkDone = errors.New("invalid work done")
	// ErrInvalidResource returned when data resource does not match client name
	ErrInvalidResource = errors.New("invalid resource")
)

type ServerPOWWorker struct {
	HashFn func() hash.Hash
	Cfg    ServerWorkerConfig
}

type ServerWorkerConfig struct {
	LeadingZeroBits uint
}

func NewServerPOWWorker(hashFn func() hash.Hash, cfg ServerWorkerConfig) *ServerPOWWorker {
	return &ServerPOWWorker{
		HashFn: hashFn,
		Cfg:    cfg,
	}
}

func (h *ServerPOWWorker) GenerateNew(_ context.Context, resource string) (pow.Data, error) {
	data := NewHashcashData(h.Cfg.LeadingZeroBits, resource)
	return data.ToPOWData(), nil
}

// ValidateWorkDone - checks that hash is valid.
// Return ErrInvalidWorkDone error is work result hash is not walid.
// Effectively this function should have more checks, for example data is not expired or Resource should be cashed
func (h *ServerPOWWorker) ValidateWorkDone(_ context.Context, resource string, data pow.Data) error {
	hashcashData, err := Parse(data)
	if err != nil {
		return err
	}

	// Just example check that beside validating hash
	// the data itself should be validated as well
	if hashcashData.Resource != resource {
		return ErrInvalidResource
	}

	// Extra check - client should increment counter at least once
	if hashcashData.Counter == 0 {
		return ErrInvalidWorkDone
	}

	hashValid, err := hashcashData.isHashValid(h.HashFn())
	if err != nil {
		return fmt.Errorf("isHashValid: %w", err)
	}

	if !hashValid {
		return ErrInvalidWorkDone
	}

	return nil
}
