package server

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
	"github.com/karpovicham/word-of-wisdom/internal/proto"
	"github.com/karpovicham/word-of-wisdom/pkg/pow"
	"github.com/karpovicham/word-of-wisdom/service/quotes_book"
)

type Server struct {
	Log       logger.Logger
	Cfg       Config
	Services  Services
	POWWorker pow.ServerWorker
	MsgrFn    func(conn net.Conn) messenger.Messenger
}

type Services struct {
	QuotesBook quotes_book.QuotesBook
}

func NewTCPServer(cfg Config, logger logger.Logger, services Services,
	powWorker pow.ServerWorker, msgrFn messenger.MsgrFn) *Server {
	return &Server{
		Cfg:       cfg,
		Log:       logger,
		Services:  services,
		POWWorker: powWorker,
		MsgrFn:    msgrFn,
	}
}

// Run - start the server and serve new client connection in endless loop
func (s *Server) Run(ctx context.Context) error {
	// Init server listener
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Cfg.Host, s.Cfg.Port))
	if err != nil {
		return err
	}
	defer ln.Close()

	s.Log.Info("Listening:", ln.Addr())

	// Track context status
	go watchContext(ctx, ln)

	// Us WG is gracefully shut down the serving process
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			if canceled(ctx) {
				return
			}

			conn, err := ln.Accept()
			if err != nil {
				switch true {
				// Error could be received for closed connection by manual interrupt
				case canceled(ctx):
					return
				default:
					s.Log.Error("Accepting conn:", err)
					time.Sleep(1 * time.Second)
					continue
				}
			}

			// Effectively this should be implemented with a workers pool
			go s.handleConnection(ctx, conn)
		}
	}()

	// Waiting for server to stop serving connections
	wg.Wait()
	return nil
}

// Context could be closed by manual interrupt.
// Close listener connection to stop serving connections and as the result - service runtime.
func watchContext(ctx context.Context, ln net.Listener) {
	<-ctx.Done()
	ln.Close()
}

// Serve client connection
func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	clientAddr := conn.RemoteAddr()
	s.Log.Info("New client:", clientAddr.String())
	defer s.Log.Info("Close client:", clientAddr.String())
	defer conn.Close()

	handler := NewRequestsHandler(s, conn)
	for {
		receivedMsg, err := handler.Msgr.Receive()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			s.Log.Error("Receive:", err)
			return
		}

		switch receivedMsg.Type {
		case proto.TypeChallenge:
			if err = handler.HandleChallengeRequest(ctx); err != nil {
				s.Log.Error("HandleChallengeRequest:", err)
				return
			}
		case proto.TypeQuote:
			if err = handler.HandleQuoteRequest(ctx, receivedMsg); err != nil {
				s.Log.Error("HandleQuoteRequest:", err)
				return
			}
		case proto.TypeStop:
			return
		default:
			s.Log.Error("Unsupported protocols:", receivedMsg.Type, err)
			return
		}
	}
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
