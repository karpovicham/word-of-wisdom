package client

import (
	"context"
	"fmt"

	"github.com/karpovicham/word-of-wisdom/service/client/resolver"
)

// ProcessQuote Makes all required steps to get the Quote in the end
func (c *Client) ProcessQuote(ctx context.Context, r resolver.Resolver) error {
	newPOWData, err := r.RequestChallenge()
	if err != nil {
		return fmt.Errorf("request challenge: %w", err)
	}

	// Worker performs some work on the data and returns the result data to be validated on the server side
	resultPOWData, err := c.POWWorker.DoWork(ctx, newPOWData)
	if err != nil {
		return fmt.Errorf("do work: %w", err)
	}

	quote, err := r.RequestQuote(resultPOWData)
	if err != nil {
		return fmt.Errorf("request quote: %w", err)
	}

	c.Log.Info("Result quote:", *quote)
	return nil
}
