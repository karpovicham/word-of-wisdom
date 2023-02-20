package server

import (
	"context"
	"fmt"
	"net"

	"github.com/karpovicham/word-of-wisdom/internal/messenger"
	"github.com/karpovicham/word-of-wisdom/internal/proto"
)

type RequestHandler struct {
	*Server
	Msgr       messenger.Messenger
	ClientName string
}

func NewRequestsHandler(server *Server, conn net.Conn) *RequestHandler {
	return &RequestHandler{
		Server:     server,
		Msgr:       server.MsgrFn(conn),
		ClientName: conn.RemoteAddr().String(),
	}
}

func (h *RequestHandler) HandleChallengeRequest(ctx context.Context) error {
	data, err := h.POWWorker.GenerateNew(ctx, h.ClientName)
	if err != nil {
		return fmt.Errorf("generateNew: %w", err)
	}

	if err := h.Msgr.Send(&proto.Message{
		Type: proto.Challenge,
		Data: data,
	}); err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func (h *RequestHandler) HandleQuoteRequest(ctx context.Context, msg *proto.Message) error {
	if err := h.POWWorker.ValidateWorkDone(ctx, h.ClientName, msg.Data); err != nil {
		return fmt.Errorf("generateNew: %w", err)
	}

	quote, err := h.Services.QuotesBook.GetRandomQuote(ctx)
	if err != nil {
		return err
	}

	if err = h.Msgr.Send(&proto.Message{
		Type: proto.Quote,
		Data: quote.ToJson(),
	}); err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}
