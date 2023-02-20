package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"os"
	"time"

	"github.com/karpovicham/word-of-wisdom/internal/logger"
	"github.com/karpovicham/word-of-wisdom/internal/signal"
	"github.com/karpovicham/word-of-wisdom/pkg/hashcash"
	"github.com/karpovicham/word-of-wisdom/service/client"
)

var (
	host           = flag.String("host", "", "Connection host")
	port           = flag.String("port", "9992", "Connection port")
	powWorkTimeout = flag.Duration("powWorkTimeout", 10*time.Second, "Timeout for the POW work to be done")
)

func main() {
	// Read cmd arguments
	flag.Parse()
	cfg := client.Config{
		Host: *host,
		Port: *port,
	}

	log := logger.NewLogger(os.Stdout)

	ctx, cancel := context.WithCancel(context.Background())
	go signal.WatchShutdown(cancel)

	powWorker := hashcash.NewClientPOWWorker(sha256.New, hashcash.ClientWorkerConfig{
		ComputeTimeout: *powWorkTimeout,
	})

	c := client.NewClient(cfg, log, powWorker)
	if err := c.Run(ctx); err != nil {
		log.Fatal("Run client:", err)
	}
}
