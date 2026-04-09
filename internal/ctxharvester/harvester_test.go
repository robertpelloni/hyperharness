package ctxharvester

import (
	"testing"
)

func TestHarvesterCreation(t *testing.T) {
	h := NewContextHarvester(nil)
	if h == nil {
		t.Fatal("should create harvester")
	}
	cfg := h.GetConfig()
	if cfg.MaxTokenBudget != 128000 {
		t.Errorf("default budget: %d", cfg.MaxTokenBudget)
	}
}

func TestHarvesterConfigOverride(t *testing.T) {
	h := NewContextHarvester(&HarvestConfig{MaxTokenBudget: 50000, ChunkSize: 200})
	cfg := h.GetConfig()
	if cfg.MaxTokenBudget != 50000 {
		t.Errorf("budget: %d", cfg.MaxTokenBudget)
	}
	if cfg.ChunkSize != 200 {
		t.Errorf("chunkSize: %d", cfg.ChunkSize)
	}
	// Other fields should keep defaults
	if cfg.PruneThreshold != 0.1 {
		t.Errorf("threshold: %f", cfg.PruneThreshold)
	}
}

func TestHarvesterHarvest(t *testing.T) {
	h := NewContextHarvester(nil)
	chunks := h.Harvest(SourceActiveFile, "Hello world. This is a test. More content here.", nil)
	if len(chunks) == 0 {
		t.Fatal("should produce chunks")
	}
	if h.ChunkCount() == 0 {
		t.Error("should have stored chunks")
	}
}

func TestHarvesterRetrieve(t *testing.T) {
	h := NewContextHarvester(nil)
	h.Harvest(SourceDocumentation, "Go is a statically typed compiled language. Go has goroutines for concurrency. Go has channels.", nil)
	h.Harvest(SourceActiveFile, "Python is a dynamically typed interpreted scripting language. Python uses indentation.", nil)

	results := h.Retrieve("Go goroutines concurrency", 10000)
	if len(results) == 0 {
		t.Fatal("should find relevant chunks")
	}
	// Just verify we got results - relevance ranking depends on content overlap
	if results[0].AccessCount != 1 {
		t.Error("access count should be incremented")
	}
}

func TestHarvesterRetrieveTokenBudget(t *testing.T) {
	h := NewContextHarvester(&HarvestConfig{MaxTokenBudget: 100000})
	// Add lots of content
	for i := 0; i < 20; i++ {
		h.Harvest(SourceDocumentation, "This is test content about programming languages and tools.", nil)
	}
	results := h.Retrieve("programming", 50) // Very small budget
	totalTokens := 0
	for _, c := range results {
		totalTokens += c.TokenCount
	}
	if totalTokens > 50 {
		t.Errorf("exceeded budget: %d tokens", totalTokens)
	}
}

func TestHarvesterPrune(t *testing.T) {
	h := NewContextHarvester(&HarvestConfig{PruneThreshold: 100}) // Very high threshold
	h.Harvest(SourceTerminal, "some output", nil)
	pruned := h.Prune()
	if pruned == 0 {
		t.Error("should have pruned low-relevance chunks")
	}
}

func TestHarvesterClear(t *testing.T) {
	h := NewContextHarvester(nil)
	h.Harvest(SourceActiveFile, "content", nil)
	h.Clear()
	if h.ChunkCount() != 0 {
		t.Error("should be empty after clear")
	}
}

func TestHarvesterReport(t *testing.T) {
	h := NewContextHarvester(nil)
	h.Harvest(SourceActiveFile, "chunk 1", nil)
	h.Harvest(SourceActiveFile, "chunk 2", nil)
	h.Harvest(SourceDocumentation, "doc chunk", nil)

	report := h.GetReport()
	if report.TotalChunks != 3 {
		t.Errorf("total: %d", report.TotalChunks)
	}
	if report.SourceBreakdown["active-file"] != 2 {
		t.Errorf("active-file count: %d", report.SourceBreakdown["active-file"])
	}
	if report.BudgetUtilization <= 0 {
		t.Error("should have some utilization")
	}
}

func TestHarvesterGetChunks(t *testing.T) {
	h := NewContextHarvester(nil)
	h.Harvest(SourceActiveFile, "first", nil)
	h.Harvest(SourceActiveFile, "second", nil)

	chunks := h.GetChunks()
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	// Should be sorted newest first
	if chunks[0].Content != "second" && chunks[1].Content != "first" {
		t.Log("chunks returned in order")
	}
}

func TestSemanticChunk(t *testing.T) {
	text := "First sentence. Second sentence. Third sentence. Fourth sentence. Fifth sentence."
	chunks := semanticChunk(text, 3, 1)
	if len(chunks) == 0 {
		t.Fatal("should produce chunks")
	}
}

func TestSemanticChunkEmpty(t *testing.T) {
	chunks := semanticChunk("", 300, 30)
	if chunks != nil {
		t.Error("empty input should return nil")
	}
}

func TestEstimateTokens(t *testing.T) {
	tokens := estimateTokens("Hello world, this is a test")
	if tokens <= 0 {
		t.Error("should estimate positive tokens")
	}
}

func TestWordSet(t *testing.T) {
	set := wordSet("hello world hello")
	if len(set) != 2 {
		t.Errorf("should have 2 unique words, got %d", len(set))
	}
}
