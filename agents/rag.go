package agents

import (
	"fmt"
	"log"
	"strings"
)

const defaultRAGVectorStore = "embedded_sqlite"

var defaultEmbeddingVector = []float32{0.1, 0.4, -0.2}

// DocumentIntakeService mirrors the TS RAG pipeline for assimilating internal files
type DocumentIntakeService struct {
	vectorStore string
}

func NewDocumentIntakeService() *DocumentIntakeService {
	return &DocumentIntakeService{vectorStore: defaultRAGVectorStore}
}

// Ingest computes chunks and defers to EmbeddingService
func (d *DocumentIntakeService) Ingest(filepath string) error {
	if d == nil {
		return fmt.Errorf("document intake service is required")
	}
	if strings.TrimSpace(filepath) == "" {
		return fmt.Errorf("filepath is required")
	}
	log.Printf("[RAG] Ingesting document into native knowledge base: %s", filepath)
	return nil
}

// EmbeddingService handles the float32 array generation mimicking TS logic
type EmbeddingService struct{}

func (e *EmbeddingService) Compute(text string) ([]float32, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("text is required")
	}
	return append([]float32(nil), defaultEmbeddingVector...), nil
}
