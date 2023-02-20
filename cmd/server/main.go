package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/karpovicham/word-of-wisdom/internal/logger"
	"github.com/karpovicham/word-of-wisdom/internal/messenger"
	"github.com/karpovicham/word-of-wisdom/internal/signal"
	"github.com/karpovicham/word-of-wisdom/pkg/hashcash"
	"github.com/karpovicham/word-of-wisdom/service/quotes_book/in_memory_book"
	"github.com/karpovicham/word-of-wisdom/service/server"
)

var (
	host            = flag.String("host", "", "Connection host")
	port            = flag.String("port", "9992", "Connection port")
	quotesFile      = flag.String("quotesFile", "assets/quotes.json", "Path to the quotes book file")
	powHashZeroBits = flag.Uint("powHashZeroBits", 20, "Number of leading zero bits required in work done hash")
)

func main() {
	// Read cmd arguments
	flag.Parse()
	cfg := server.Config{
		Host: *host,
		Port: *port,
	}

	log := logger.NewLogger(os.Stdout)

	// This should probably be a separate service
	quotesBookService, err := in_memory_book.NewQuotesBookService(
		*quotesFile,
		rand.New(rand.NewSource(time.Now().UnixNano())),
	)
	if err != nil {
		log.Fatal("Init quotes book service:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go signal.WatchShutdown(cancel)

	services := server.Services{
		QuotesBook: quotesBookService,
	}

	powWorker := hashcash.NewServerPOWWorker(sha256.New, hashcash.ServerWorkerConfig{
		LeadingZeroBits: *powHashZeroBits,
	})

	s := server.NewTCPServer(cfg, log, services, powWorker, messenger.NewMessenger)
	if err := s.Run(ctx); err != nil {
		log.Fatal("Run server:", err)
	}
}
