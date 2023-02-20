package hashcash

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"log"
	"testing"
	"time"
)

func Benchmark_HashcashComputeData(b *testing.B) {
	ctx := context.Background()
	appName := "BenchmarkTest"
	timeout := 30 * time.Second

	bench := func(b *testing.B, hashFunc func() hash.Hash, bits uint) {
		data := NewHashcashData(bits, appName)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := data.ComputeData(ctx, hashFunc(), timeout)
			if err != nil {
				log.Println("Error:", err)
				return
			}
		}
	}

	bits := []uint{
		15,
		20, // Default
	}

	b.Run("Hasher sha1", func(b *testing.B) {
		for _, value := range bits {
			b.Run(fmt.Sprintf("bits=%d", value), func(b *testing.B) {
				bench(b, sha1.New, value)
			})
		}
	})

	b.Run("Hasher sha256", func(b *testing.B) {
		for _, value := range bits {
			b.Run(fmt.Sprintf("bits=%d", value), func(b *testing.B) {
				bench(b, sha256.New, value)
			})
		}
	})
}
