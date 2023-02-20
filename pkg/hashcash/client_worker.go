package hashcash

import (
	"context"
	"hash"
	"time"

	"github.com/karpovicham/word-of-wisdom/pkg/pow"
)

type ClientPOWWorker struct {
	HashFn func() hash.Hash
	Cfg    ClientWorkerConfig
}

type ClientWorkerConfig struct {
	ComputeTimeout time.Duration
}

func NewClientPOWWorker(hashFn func() hash.Hash, cfg ClientWorkerConfig) *ClientPOWWorker {
	return &ClientPOWWorker{
		HashFn: hashFn,
		Cfg:    cfg,
	}
}

func (h *ClientPOWWorker) DoWork(ctx context.Context, data pow.Data) (pow.Data, error) {
	hashcashData, err := Parse(data)
	if err != nil {
		return nil, err
	}

	newData, err := hashcashData.ComputeData(ctx, h.HashFn(), h.Cfg.ComputeTimeout)
	if err != nil {
		return nil, err
	}

	return newData.ToPOWData(), nil
}
