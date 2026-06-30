package bench

import (
	"testing"

	"github.com/robertpelloni/hyperharness/tools"
)

func BenchmarkToolInvocationFoundation(b *testing.B) {
	registry := tools.NewRegistry()
	tool, ok := registry.Find("read")
	if !ok {
		b.Fatal("Failed to find tool")
	}

	args := map[string]interface{}{
		"path": "bench_test_file.txt",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tool.Execute(args)
	}
}

func BenchmarkToolInvocationParity(b *testing.B) {
	registry := tools.NewRegistry()
	tool, ok := registry.Find("read_file") // gemini opencode parity equivalent to read
	if !ok {
		b.Fatal("Failed to find tool")
	}

	args := map[string]interface{}{
		"file_path": "bench_test_file.txt",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tool.Execute(args)
	}
}
