package client

import (
	"context"
	"fmt"

	"github.com/karpovicham/word-of-wisdom/service/client/resolver"
)

func (c *Client) ProcessQuote(ctx context.Context, r resolver.Resolver) error {
	powData, err := r.RequestChallenge()
	if err != nil {
		return fmt.Errorf("request challenge: %w", err)
	}

	newPOWData, err := c.POWWorker.DoWork(ctx, powData)
	if err != nil {
		return fmt.Errorf("DoWork: %w", err)
	}

	quote, err := r.RequestQuote(newPOWData)
	if err != nil {
		return fmt.Errorf("request quote: %w", err)
	}

	c.Log.Info("Result quote:", *quote)
	return nil
}
