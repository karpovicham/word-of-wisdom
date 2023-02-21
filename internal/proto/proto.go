// Package proto - describes messages structure for Client-Server TCP connection
package proto

import (
	"encoding/json"
	"errors"
)

// Type represents request/response type
// Client/Server apps should detect how to handle requests by this value
type Type int

const (
	// TypeChallenge - Client challenge request, server challenge response
	TypeChallenge Type = iota
	// TypeQuote - Client challenge resolve request, server quote response
	TypeQuote
	// TypeStop - Message to close connection for both sides
	TypeStop
)

type Error error

var (
	// ErrorInvalidData - Request data is not valid (ie not correct data format for the Type)
	ErrorInvalidData = Error(errors.New("invalid data"))
	// ErrorNotVerified - Verification failed or Verification required to access the resource
	ErrorNotVerified = Error(errors.New("not verified"))
)

// Message - Transferred data structure
type Message struct {
	Type  Type   `json:"type"`
	Data  []byte `json:"data"`
	Error *Error `json:"error,omitempty"`
}

// Parse - decodes JSON to Message
func Parse(b []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// ToString - encodes Message to JSON string
func (m *Message) ToString() string {
	msgBytes, _ := json.Marshal(m)
	return string(msgBytes)
}

// ToJSON - encodes Message to JSON bytes
func (m *Message) ToJSON() []byte {
	msgBytes, _ := json.Marshal(m)
	return msgBytes
}
