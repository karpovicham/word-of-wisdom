// Package signal provides helper method for graceful shutdown of applications
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WatchShutdown - subscribes to Interrupt like signals and stops context
func WatchShutdown(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)

	// Register signals that stops application (plus USR1 one for testing)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1)

	// Wait for a signal
	<-sig

	// Got signal, cancel context
	cancel()
}
