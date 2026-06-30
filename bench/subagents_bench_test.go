package bench

import (
	"testing"

	"github.com/robertpelloni/hyperharness/internal/subagents"
)

func BenchmarkSubagentSpawning(b *testing.B) {
	manager := subagents.NewManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.CreateTask(subagents.TypeCode, "Bench prompt", "Bench input", "Bench context")
	}
}
