package bench

import (
	"fmt"
	"strings"
	"testing"

	"github.com/robertpelloni/hyperharness/internal/context"
)

func BenchmarkContextCompaction(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				mgr := context.NewManager(1000, 100000)
				mgr.EnableAutoCompact(true)
				for j := 0; j < size; j++ {
					mgr.Add(context.Message{
						Role:    "user",
						Content: "This is a simulated message content for compaction benchmark. " + strings.Repeat("A ", 10),
					})
				}
				b.StartTimer()
				mgr.Compact(size / 2)
			}
		})
	}
}
