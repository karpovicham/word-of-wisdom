// Package proto - describes messages structure for Client-Server TCP connection
package proto

import "encoding/json"

// Type represents request/response type
// Client/Server apps should detect how to handle requests by this value
type Type int

const (
	// Challenge - Client challenge request, server challenge response
	Challenge Type = iota
	// Quote - Client challenge resolve request, server quote response
	Quote
	// Stop - Message to close connection for both sides
	Stop
)

// Message - Transferred data structure
type Message struct {
	Type Type   `json:"type"`
	Data []byte `json:"data"`
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
