// Package pow contains interfaces and data types for POW workers with different algorithms
package pow

import "context"

// Data represents data type for POW workers
type Data []byte

// ClientWorker - interface for POW worker on the client side
type ClientWorker interface {
	// DoWork perform work on the Data and return modified Data to Validate on Server side
	DoWork(ctx context.Context, data Data) (Data, error)
}

// ServerWorker - interface for POW worker on the server side
type ServerWorker interface {
	// GenerateNew generates new work Data
	GenerateNew(ctx context.Context, clientName string) (Data, error)
	// ValidateWorkDone check if the work is really done and the Data is valid
	ValidateWorkDone(ctx context.Context, clientName string, data Data) error
}
