package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/karpovicham/word-of-wisdom/internal/logger"
	"github.com/karpovicham/word-of-wisdom/internal/messenger"
	"github.com/karpovicham/word-of-wisdom/pkg/pow"
	"github.com/karpovicham/word-of-wisdom/service/client/resolver"
)

type Client struct {
	Log       logger.Logger
	Cfg       Config
	POWWorker pow.ClientWorker
}

func NewClient(cfg Config, log logger.Logger, powWorker pow.ClientWorker) *Client {
	return &Client{
		Cfg:       cfg,
		Log:       log,
		POWWorker: powWorker,
	}
}

// Run connects to the server
// and start getting Quote every 3 seconds in endless loop
func (c *Client) Run(ctx context.Context) error {
	// Create a socket and start listening
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", c.Cfg.Host, c.Cfg.Port))
	if err != nil {
		return fmt.Errorf("deal: %w", err)
	}
	defer conn.Close()

	c.Log.Info("Client", conn.LocalAddr(), "connected to:", conn.RemoteAddr())

	// Track context status
	go watchContext(ctx, conn)

	apiResolver := resolver.NewClientAPIResolver(messenger.NewMessenger(conn))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Endless loop to get a new Quote every 3 secs
		for {
			if err = c.ProcessQuote(ctx, apiResolver); err != nil {
				switch true {
				case canceled(ctx):
				case errors.Is(err, io.EOF):
					c.Log.Error("Server closed connection")
				default:
					c.Log.Error("Process quote: ", err)
				}
				return
			}

			time.Sleep(3 * time.Second)
		}
	}()

	wg.Wait()
	return nil
}

// Context could be closed by manual interrupt.
// Close listener connection to stop serving connections and as the result - service runtime.
func watchContext(ctx context.Context, conn net.Conn) {
	<-ctx.Done()
	conn.Close()
}

// canceled returns true if the context is canceled
func canceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
