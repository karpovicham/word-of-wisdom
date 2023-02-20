package hashcash

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"math/rand"
	"strconv"
	"time"

	"github.com/karpovicham/word-of-wisdom/pkg/pow"
)

const (
	headerFormat               = "%d:%d:%s:%s::%s:%s"
	hashcashFormatVersion uint = 1
)

var (
	ErrComputeTimeout = errors.New("timeout")
	ErrClosedContext  = errors.New("closed context")
)

// data contains data to wort with Hashcash algorithm
// ver: Hashcash format Version, 1 (which supersedes version 0).
// bits: Number of "partial pre-image" (zero) bits in the hashed code.
// date: The time that the message was sent, in the format YYMMDD[hhmm[ss]] (changed to Unix Micro timestamp)
// resource: Resource data string being transmitted, e.g., an IP address or email address.
// ext: Extension (optional; ignored in version 1).
// rand: String of random characters, encoded in base-64 format.
// counter: Binary counter, encoded in base-64 format.
// https://en.wikipedia.org/wiki/Hashcash
type data struct {
	Version  uint   `json:"version"`
	Bits     uint   `json:"bits"`
	Date     string `json:"date"`
	Resource string `json:"resource"`
	Random   string `json:"random"`
	Counter  int64  `json:"counter"`
}

// NewHashcashData returns newly generated ready to work Hashcash data
func NewHashcashData(leadingZeroBits uint, resource string) *data {
	return &data{
		Version:  hashcashFormatVersion,
		Bits:     leadingZeroBits,
		Date:     strconv.FormatInt(time.Now().UnixMicro(), 10),
		Resource: resource,
		Random:   base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(int64(rand.Int()), 10))),
		Counter:  0, // This should be incremented on client side
	}
}

// Header - Formats the Hashcash data to known Header format X-Hashcash
// Example: 1:20:1303030600:anni@cypherspace.org::McMybZIhxKXu57jd:ckvi
func (d *data) Header() string {
	counter := base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(d.Counter, 10)))
	return fmt.Sprintf(headerFormat, hashcashFormatVersion, d.Bits, d.Date, d.Resource, d.Random, counter)
}

// ComputeData - hash data till it meets the zero leading Bits criteria
// and return new data with incremented counter
func (d data) ComputeData(ctx context.Context, hasher hash.Hash, computeTimeout time.Duration) (*data, error) {
	// Validate Bits
	if int(d.Bits) > hasher.Size()*8 {
		return nil, fmt.Errorf("bits config %d is out of the hash range (%d bits)", d.Bits, hasher.Size()*8)
	}

	// Loop should be controlled by timeout ot iterations count
	timer := time.NewTimer(computeTimeout)
	defer timer.Stop()

	// Hash header of hashcash data till it meets leading zero bits criteria
	for {
		// Ways to escape from endless loop
		// Effectively iterations limit should be implemented as well
		select {
		case <-ctx.Done():
			return nil, ErrClosedContext
		case <-timer.C:
			return nil, ErrComputeTimeout
		default:
		}

		d.Counter++

		headerHash := hashData(d.Header(), hasher)
		if leadingZeroBitsMatched(headerHash, d.Bits) {
			return &d, nil
		}
	}
}

// Valid - checks that hash meets the zero leading Bits criteria
func (d *data) isHashValid(hasher hash.Hash) (bool, error) {
	// Validate Bits
	if int(d.Bits) > hasher.Size()*8 {
		return false, fmt.Errorf("bits config %d is out of the hash range (%d bits)", d.Bits, hasher.Size()*8)
	}

	// Check that
	headerHash := hashData(d.Header(), hasher)
	if !leadingZeroBitsMatched(headerHash, d.Bits) {
		return false, nil
	}

	return true, nil
}

// Parse - decodes JSON to Hashcash data
func Parse(b []byte) (*data, error) {
	var data data
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// ToPOWData - encodes data to JSON bytes
func (d *data) ToPOWData() pow.Data {
	bytes, _ := json.Marshal(d)
	return bytes
}

// hashData returns hash of the in string
func hashData(in string, h hash.Hash) []byte {
	h.Reset()
	h.Write([]byte(in))
	return h.Sum(nil)
}

// leadingZeroBitsMatched return true if hash first N bits are zero
func leadingZeroBitsMatched(hash []byte, bits uint) bool {
	// Check quotient bytes are zero
	quotient := bits / 8
	for _, b := range hash[:quotient] {
		if b != 0 {
			return false
		}
	}

	modulo := bits % 8

	// All bits are checked as whole bytes
	if modulo == 0 {
		return true
	}

	// Check the left bits in the next byte
	b := hash[quotient]
	switch modulo {
	case 1:
		if b > 127 {
			return false
		}
	case 2:
		if b > 63 {
			return false
		}
	case 3:
		if b > 31 {
			return false
		}
	case 4:
		if b > 15 {
			return false
		}
	case 5:
		if b > 7 {
			return false
		}
	case 6:
		if b > 3 {
			return false
		}
	case 7:
		if b > 1 {
			return false
		}
	}

	return true
}
